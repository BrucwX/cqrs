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
	Create(ctx context.Context, t *teacher.Teacher) error
	// Update 更新讲师：按 ID 取出已有讲师交给 updateFn 改，改完写回
	//
	// 讲师不存在时报 ErrTeacherNotFound。
	Update(
		ctx context.Context,
		id int64,
		updateFn func(ctx context.Context, t *teacher.Teacher) (*teacher.Teacher, error),
	) error
	// Delete 删除讲师
	Delete(ctx context.Context, id int64) error
	// Get 取讲师；不存在时返回 (nil, nil)
	Get(ctx context.Context, id int64) (*teacher.Teacher, error)
	// MustGet 取讲师聚合本身；取不到报 ErrTeacherNotFound
	//
	// 与 Get 的差别：Get 找不到时是 (nil, nil)（给「查到了没」的调用方），
	// MustGet 是规则判定要用的，找不到必须报错。
	MustGet(ctx context.Context, teacherID int64) (teacher.Teacher, error)
}
