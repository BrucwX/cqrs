package command

import (
	"context"
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
	// Create 新增学员
	Create(ctx context.Context, s *student.Student) error
	// Update 更新学员：按 ID 取出已有学员交给 updateFn 改，改完写回
	//
	// 学员不存在时报 ErrStudentNotFound。
	Update(
		ctx context.Context,
		id int64,
		updateFn func(ctx context.Context, s *student.Student) (*student.Student, error),
	) error
	// Delete 删除学员
	Delete(ctx context.Context, id int64) error
	// Get 取学员；不存在时返回 (nil, nil)
	Get(ctx context.Context, id int64) (*student.Student, error)
}
