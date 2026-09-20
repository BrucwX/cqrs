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
	// AfterFunc 在延迟 d 之后在独立协程中执行任务 f
	AfterFunc(d time.Duration, task Task) (Timer, error)

	// Every 周期性执行任务，interval 为执行间隔
	Every(interval time.Duration, task Task) (Timer, error)

	// Close 优雅停止调度器，释放底层资源
	Close(ctx context.Context) error
}