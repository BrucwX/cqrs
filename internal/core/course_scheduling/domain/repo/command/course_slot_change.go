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
	Save(csc *courseSlotChange.CourseSlotChange) error
	// Delete 删除课表变更
	Delete(id int64) error
	// Change 登记一次临时换课
	//
	// checkConflictFn 由调用方注入，仓库会把「本次换课」传进去；
	// 传 nil 表示不做检查。冲突时不写入。
	Change(ctx context.Context, csc *courseSlotChange.CourseSlotChange, checkConflictFn func(ctx context.Context, csc *courseSlotChange.CourseSlotChange) (bool, error)) error
}
