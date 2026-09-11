package command

import (
	"errors"

	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
)

var (
	// ErrTeacherRequired 传入的讲师为空。
	ErrTeacherRequired = errors.New("teacher is required")
	// ErrTeacherNotFound 指定的讲师不存在。
	ErrTeacherNotFound = errors.New("teacher not found")
)

// TeacherCommand 讲师命令接口
type TeacherCommand interface {
	// Save 保存讲师（新增或更新）
	Save(t *teacher.Teacher) error
	// Delete 删除讲师
	Delete(id int64) error
}
