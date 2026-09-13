package command

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
)

// TeacherCommand 讲师命令接口
type TeacherCommand interface {
	// CreateTeacher 新增讲师
	CreateTeacher(ctx context.Context, t *teacher.Teacher) error
	// UpdateTeacher 更新讲师：按 ID 取出已有讲师交给 updateFn 改，改完写回
	//
	// 讲师不存在时报 ErrTeacherNotFound。
	UpdateTeacher(
		ctx context.Context,
		id int64,
		updateFn func(ctx context.Context, t *teacher.Teacher) (*teacher.Teacher, error),
	) error
	// DeleteTeacher 删除讲师
	DeleteTeacher(ctx context.Context, id int64) error
	// GetTeacher 取讲师；不存在时返回 (nil, nil)
	GetTeacher(ctx context.Context, id int64) (*teacher.Teacher, error)
	// MustGetTeacher 取讲师聚合本身；取不到报 ErrTeacherNotFound
	//
	// 与 Get 的差别：Get 找不到时是 (nil, nil)（给「查到了没」的调用方），
	// MustGetTeacher 是规则判定要用的，找不到必须报错。
	MustGetTeacher(ctx context.Context, teacherID int64) (teacher.Teacher, error)
}
