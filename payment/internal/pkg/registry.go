package timer

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"
)

var (
	// ErrEmptyJobName 任务必须有名字（注册 / 停止 / 手动触发都靠它指认）。
	ErrEmptyJobName = errors.New("timer: job name is empty")
	// ErrDuplicateJob 同名任务已经注册过了。
	ErrDuplicateJob = errors.New("timer: job already registered")
	// ErrJobNotFound 没有这个名字的任务。
	ErrJobNotFound = errors.New("timer: job not found")
	// ErrInvalidSpec 调度规格不合法：Delay 与 Every 必须且只能设一个。
	ErrInvalidSpec = errors.New("timer: spec must set exactly one of Delay / Every")
)

// JobName 任务名：全局唯一。
type JobName string

// Spec 调度规格：**二选一** —— 延迟跑一次（Delay）或周期跑（Every）。
type Spec struct {
	Delay time.Duration // 只跑一次：延迟多久后跑
	Every time.Duration // 周期跑：每隔多久跑
}

// Job 一个可注册的定时任务：**任务自己说清楚叫什么、怎么调度、干什么**。
//
// 任务体只关心业务（扫一批过期单、对账……），调度多久由 Spec 声明，注册与启动交给 Registry。
type Job interface {
	Name() JobName                 // 唯一名
	Spec() Spec                    // 怎么调度
	Run(ctx context.Context) error // 任务体
}

// Registry 任务注册器：把「有哪些定时任务」从调度器里拎出来集中登记。
//
// 分工：
//   - Registry 管**有哪些任务**（注册、按名启停、手动触发、查列表）；
//   - Scheduler 管**怎么跑**（延迟 / 周期、协程、句柄）。
//
// 进程启动时把各处的 Job 注册进来，Start 一次性挂到 Scheduler 上（任务 id 就用 Job.Name()）。
//
// ⚠️ 常驻任务（周期扫描、对账）走注册器；像「这张单 30 分钟后关掉」这种**一次性的动态**定时，
// 仍可直接用 Scheduler.AfterFunc —— 不过更稳的写法是注册一个周期任务去扫「已过期的待支付订单」
// （重启不丢、多实例天然安全）。
type Registry interface {
	Register(jobs ...Job) error                      // 注册（空名 / 重名报错）
	Start(ctx context.Context) error                 // 启动所有已注册任务
	Stop(names ...JobName) error                     // 停指定任务（不传 = 停全部）
	Trigger(ctx context.Context, name JobName) error // 手动跑一次（运维 / 补跑 / 测试）
	Names() []JobName                                // 已注册的任务名（升序）
}

// NewRegistry 构造注册器：需要一个底层调度器。
func NewRegistry(scheduler Scheduler) Registry {
	return &registry{
		scheduler: scheduler,
		jobs:      make(map[JobName]Job),
		live:      make(map[JobName]Timer),
	}
}

// registry Registry 的默认实现。
//
// 注册表与运行中的句柄都加锁保护：注册通常在启动阶段，触发 / 停止可能来自运维接口。
type registry struct {
	scheduler Scheduler

	mu   sync.Mutex
	jobs map[JobName]Job   // 已注册的任务
	live map[JobName]Timer // Start 之后正在跑的任务句柄
}

// Register 注册一批任务：名字必须非空且不重复。
func (r *registry) Register(jobs ...Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, j := range jobs {
		name := j.Name()
		if name == "" {
			return ErrEmptyJobName
		}
		if _, ok := r.jobs[name]; ok {
			return fmt.Errorf("%w: %s", ErrDuplicateJob, name)
		}
		r.jobs[name] = j
	}
	return nil
}

// Start 把已注册的任务挂到调度器上：Every 走周期，Delay 走一次性。
func (r *registry) Start(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for name, j := range r.jobs {
		if _, ok := r.live[name]; ok {
			continue // 已经在跑了，别重复挂
		}

		spec := j.Spec()
		t, err := r.schedule(name, spec, j.Run)
		if err != nil {
			return err
		}
		r.live[name] = t
	}
	return nil
}

// Stop 停掉指定任务（不传名字 = 停全部）。
func (r *registry) Stop(names ...JobName) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(names) == 0 {
		for name, t := range r.live {
			t.Stop()
			delete(r.live, name)
		}
		return nil
	}

	for _, name := range names {
		t, ok := r.live[name]
		if !ok {
			return fmt.Errorf("%w: %s", ErrJobNotFound, name)
		}
		t.Stop()
		delete(r.live, name)
	}
	return nil
}

// Trigger 手动跑一次指定任务（同步执行，返回任务体的错误）。
func (r *registry) Trigger(ctx context.Context, name JobName) error {
	r.mu.Lock()
	j, ok := r.jobs[name]
	r.mu.Unlock()

	if !ok {
		return fmt.Errorf("%w: %s", ErrJobNotFound, name)
	}
	return j.Run(ctx)
}

// Names 列出已注册的任务名（升序，方便展示与测试断言）。
func (r *registry) Names() []JobName {
	r.mu.Lock()
	defer r.mu.Unlock()

	return slices.Sorted(func(yield func(JobName) bool) {
		for name := range r.jobs {
			if !yield(name) {
				return
			}
		}
	})
}

// schedule 按规格把任务挂到调度器上：id 用任务名，方便之后 Cancel。
func (r *registry) schedule(name JobName, spec Spec, task Task) (Timer, error) {
	switch {
	case spec.Every > 0 && spec.Delay > 0:
		return nil, fmt.Errorf("%w: %s", ErrInvalidSpec, name)
	case spec.Every > 0:
		return r.scheduler.Every(string(name), spec.Every, task)
	case spec.Delay > 0:
		return r.scheduler.AfterFunc(string(name), spec.Delay, task)
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidSpec, name)
	}
}
