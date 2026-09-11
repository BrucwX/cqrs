package command

import (
	"context"
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
	// Create 新增讲师
	Create(t *teacher.Teacher) error
	// Update 更新讲师：按 ID 取出已有讲师交给 updateFn 改，改完写回
	//
	// 讲师不存在时报 ErrTeacherNotFound。
	Update(
		ctx context.Context,
		id int64,
		updateFn func(ctx context.Context, t *teacher.Teacher) (*teacher.Teacher, error),
	) error
	// Delete 删除讲师
	Delete(id int64) error
	// Get 取讲师；不存在时返回 (nil, nil)
	Get(id int64) (*teacher.Teacher, error)
}
