package command

import (
	"context"
	"errors"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
)

var (
	// ErrMakeupRequired 传入的补课预约为空。
	ErrMakeupRequired = errors.New("student makeup is required")
	// ErrMakeupNotFound 指定的补课预约不存在。
	ErrMakeupNotFound = errors.New("student makeup not found")
	// ErrMakeupConflict 补课预约被拒绝（目标那节课的教室装不下）。
	ErrMakeupConflict = errors.New("student makeup conflicts with classroom capacity")
)

// StudentMakeupCommand 补课申请命令接口
type StudentMakeupCommand interface {
	// Save 保存补课申请（新增或更新）
	Save(ctx context.Context, m *makeup.StudentMakeup) error
	// Delete 删除补课申请
	Delete(ctx context.Context, id int64) error
	// MustGet 取补课预约聚合；不存在时报 ErrMakeupNotFound
	MustGet(ctx context.Context, id int64) (makeup.StudentMakeup, error)
	// GetMakeupsForTarget 取补到同一节课上的全部补课预约
	//
	// 「同一节课」按槽位 + 日期两把钥匙认：同一个周排槽位在别的日期上是另一堂
	// 课，不抢这间教室的座位。状态不过滤，要不要算上由调用方决定。
	GetMakeupsForTarget(ctx context.Context, targetSlotID string, targetDate time.Time) ([]makeup.StudentMakeup, error)
}
