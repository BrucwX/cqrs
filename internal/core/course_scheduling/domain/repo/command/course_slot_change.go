package command

import (
	"context"
	"errors"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
)

var (
	// ErrSlotChangeRequired 传入的课表变更单为空。
	ErrSlotChangeRequired = errors.New("course slot change is required")
	// ErrSlotChangeNotFound 指定的课表变更单不存在。
	ErrSlotChangeNotFound = errors.New("course slot change not found")
	// ErrSlotChangeConflict 换课被拒绝（目标讲师或教室在目标时段已被占用）。
	ErrSlotChangeConflict = errors.New("course slot change conflicts with existing schedule")
)

// CourseSlotChangeCommand 课表变更命令接口
type CourseSlotChangeCommand interface {
	// Save 保存课表变更（新增或更新，不做检查）
	Save(ctx context.Context, csc *courseSlotChange.CourseSlotChange) error
	// Delete 删除课表变更
	Delete(ctx context.Context, id int64) error
	// Change 登记一次临时换课
	//
	// 只管写：把这张变更单落库。冲突判定（目标讲师 / 教室在目标时段是否已被占用）
	// 不在这里 —— 那是调用方的事（见 domain/service/schedule 的 Conflict.CheckSlotChange），
	// 判定通过才调过来。所以这里没有回调，也没有「传 nil 表示不检查」这类分支。
	Change(ctx context.Context, csc *courseSlotChange.CourseSlotChange) error
	// todo 这里不应该取这个，应该取 和 CourseSlotChange 有交集的换课记录
	// GetOtherSlotChanges 取除该变更单以外的全部换课记录
	GetOtherSlotChanges(ctx context.Context, id int64) ([]*courseSlotChange.CourseSlotChange, error)
}
