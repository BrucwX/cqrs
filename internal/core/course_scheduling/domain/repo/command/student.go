package command

import (
	"errors"

	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
)

var (
	// ErrStudentRequired 传入的学员为空。
	ErrStudentRequired = errors.New("student is required")
	// ErrStudentNotFound 指定的学员不存在。
	ErrStudentNotFound = errors.New("student not found")
)

// StudentCommand 学员命令接口
type StudentCommand interface {
	// Save 保存学员（新增或更新）
	Save(s *student.Student) error
	// Delete 删除学员
	Delete(id int64) error
}
