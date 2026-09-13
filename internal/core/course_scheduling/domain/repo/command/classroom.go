package command

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
)

// ClassroomCommand 教室命令接口
type ClassroomCommand interface {
	// CreateClassroom 新增教室
	CreateClassroom(ctx context.Context, c *classroom.Classroom) error
	// UpdateClassroom 更新教室：按 ID 取出已有教室交给 updateFn 改，改完写回
	//
	// 教室不存在时报 ErrClassroomNotFound。
	UpdateClassroom(
		ctx context.Context,
		id string,
		updateFn func(ctx context.Context, cl *classroom.Classroom) (*classroom.Classroom, error),
	) error
	// DeleteClassroom 删除教室
	DeleteClassroom(ctx context.Context, id string) error
	// MustGetClassroom 取教室聚合本身；不存在时返回 ErrClassroomNotFound
	MustGetClassroom(ctx context.Context, id string) (classroom.Classroom, error)
}
