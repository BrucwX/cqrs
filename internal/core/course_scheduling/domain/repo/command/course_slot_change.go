package command

import "cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"

// CourseSlotChangeCommand 课表变更命令接口
type CourseSlotChangeCommand interface {
	// Save 保存课表变更（新增或更新）
	Save(csc *courseSlotChange.CourseSlotChange) error
	// Delete 删除课表变更
	Delete(id int64) error
}
