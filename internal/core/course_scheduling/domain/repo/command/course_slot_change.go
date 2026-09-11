package command

import (
	"context"
	"errors"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
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

// SlotChangeContext 是仓库侧在换课时一并交给冲突检查的上下文。
//
// 仓库（命令适配器）能直接读到读模型，所以由它把 checkSlotChange 需要的
// 「本次换课 / 目标讲师的周排期 / 目标教室的周排期 / 其他换课记录」装好传进来。
type SlotChangeContext struct {
	// Change 本次换课，含目标具体时间段与目标讲师、教室
	Change courseSlotChange.CourseSlotChange
	// TeacherSlots 目标讲师的周排期（已排除本次换课归属课程自身的排期，
	// 因为那就是被换掉的那节课）
	TeacherSlots courseSlot.CourseSlots
	// ClassroomSlots 目标教室的周排期（同样排除本次换课归属课程自身的排期）
	ClassroomSlots courseSlot.CourseSlots
	// OtherChanges 其他换课记录（已排除本次换课本身）
	OtherChanges []*courseSlotChange.CourseSlotChange
}

// CourseSlotChangeCommand 课表变更命令接口
type CourseSlotChangeCommand interface {
	// Save 保存课表变更（新增或更新，不做检查）
	Save(csc *courseSlotChange.CourseSlotChange) error
	// Delete 删除课表变更
	Delete(id int64) error
	// Change 登记一次临时换课
	//
	// checkConflictFn 由调用方注入，仓库会把 SlotChangeContext 装好传进去；
	// 传 nil 表示不做检查。冲突时不写入。
	Change(ctx context.Context, csc *courseSlotChange.CourseSlotChange, checkConflictFn func(ctx context.Context, sc SlotChangeContext) (bool, error)) error
}
