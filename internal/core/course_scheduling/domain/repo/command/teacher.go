package command

import "cqrs/internal/core/course_scheduling/domain/aggregate/teacher"

// TeacherCommand 讲师命令接口
type TeacherCommand interface {
	// Save 保存讲师（新增或更新）
	Save(t *teacher.Teacher) error
	// Delete 删除讲师
	Delete(id int64) error
}
