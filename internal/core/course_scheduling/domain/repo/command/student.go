package command

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
)

// StudentCommand 学员命令接口
type StudentCommand interface {
	// CreateStudent 新增学员
	CreateStudent(ctx context.Context, s *student.Student) error
	// UpdateStudent 更新学员：按 ID 取出已有学员交给 updateFn 改，改完写回
	//
	// 学员不存在时报 ErrStudentNotFound。
	UpdateStudent(
		ctx context.Context,
		id int64,
		updateFn func(ctx context.Context, s *student.Student) (*student.Student, error),
	) error
	// DeleteStudent 删除学员
	DeleteStudent(ctx context.Context, id int64) error
	// GetStudent 取学员；不存在时返回 (nil, nil)
	GetStudent(ctx context.Context, id int64) (*student.Student, error)
	// MustGetStudent 取学员聚合本身；取不到报 ErrStudentNotFound
	MustGetStudent(ctx context.Context, id int64) (student.Student, error)
}
