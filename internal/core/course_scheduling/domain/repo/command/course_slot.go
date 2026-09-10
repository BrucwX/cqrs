package command

import "cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"

// CourseSlotCommand 课表槽位命令接口
type CourseSlotCommand interface {
	// Save 保存课表槽位（新增或更新）
	Save(cs *courseSlot.CourseSlot) error
	// Delete 删除课表槽位
	Delete(id int64) error
}
