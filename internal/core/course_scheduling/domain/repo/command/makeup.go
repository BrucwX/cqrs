package command

import (
	"context"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
)

// StudentMakeupCommand 补课申请命令接口
type StudentMakeupCommand interface {
	// SaveMakeup 保存补课申请（新增或更新）
	SaveMakeup(ctx context.Context, m *makeup.StudentMakeup) error
	// DeleteMakeup 删除补课申请
	DeleteMakeup(ctx context.Context, id int64) error
	// MustGetMakeup 取补课预约聚合；不存在时报 ErrMakeupNotFound
	MustGetMakeup(ctx context.Context, id int64) (makeup.StudentMakeup, error)
	// GetMakeupsForTarget 取补到同一节课上的全部补课预约
	//
	// 「同一节课」按槽位 + 日期两把钥匙认：同一个周排槽位在别的日期上是另一堂
	// 课，不抢这间教室的座位。状态不过滤，要不要算上由调用方决定。
	GetMakeupsForTarget(ctx context.Context, targetSlotID string, targetDate time.Time) ([]makeup.StudentMakeup, error)
}
