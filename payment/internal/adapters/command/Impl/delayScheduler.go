package command

import (
	"context"
	"time"

	"payment/internal/ports"
)

// Notify 生产延时任务，转发给 RedisData。
func (d *CommandImpl) Notify(ctx context.Context, productType string, orderID int64, delay time.Duration) error {
	return d.RedisData.Notify(ctx, productType, orderID, delay)
}

// StartCloseCourse 跑课程关单的消费循环，转发给 RedisData（阻塞到 ctx 结束）。
func (d *CommandImpl) StartCloseCourse(ctx context.Context, h ports.CourseHandler) {
	d.RedisData.StartCloseCourse(ctx, h)
}

// 编译期约束：CommandImpl 必须实现 ports.DelayScheduler。
var _ ports.DelayScheduler = (*CommandImpl)(nil)
