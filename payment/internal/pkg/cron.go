package timer

import (
	"context"
	"time"
)

// Task 定义定时器触发时执行的业务逻辑
type Task func(ctx context.Context) error

// Timer 表示单个定时器句柄
type Timer interface {
	// C 返回接收到期事件的只读通道（对齐标准库惯用语法）
	C() <-chan time.Time

	// Reset 重新设置定时器到期时间，返回重置前是否处于激活状态
	Reset(d time.Duration) bool

	// Stop 停止定时器，防止其触发，返回停止前是否处于激活状态
	Stop() bool
}

// Scheduler 负责管理和调度所有定时任务
type Scheduler interface {
	// AfterFunc 注册一个「延迟 d 之后执行一次」的任务，返回它的句柄。
	// id 全局唯一，之后可用 Cancel(id) 精确停掉这一个。
	AfterFunc(id string, d time.Duration, task Task) (Timer, error)

	// Every 注册一个「每隔 interval 执行一次」的周期任务，返回它的句柄。
	Every(id string, interval time.Duration, task Task) (Timer, error)

	// Cancel 停止并移除指定 id 的任务；返回之前是否存在（false = 没这个任务 / 已停）。
	Cancel(id string) bool

	// Close 优雅停止调度器，释放底层资源（等价于停掉所有任务）。
	Close(ctx context.Context) error
}
