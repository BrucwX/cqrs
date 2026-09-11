package command

import (
	"context"
	"errors"

	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
)

var (
	// ErrClassroomRequired 传入的教室为空。
	ErrClassroomRequired = errors.New("classroom is required")
	// ErrClassroomNotFound 指定的教室不存在。
	ErrClassroomNotFound = errors.New("classroom not found")
)

// ClassroomCommand 教室命令接口
type ClassroomCommand interface {
	// Create 新增教室
	Create(c *classroom.Classroom) error
	// Update 更新教室：按 ID 取出已有教室交给 updateFn 改，改完写回
	//
	// 教室不存在时报 ErrClassroomNotFound。
	Update(
		ctx context.Context,
		id string,
		updateFn func(ctx context.Context, cl *classroom.Classroom) (*classroom.Classroom, error),
	) error
	// Delete 删除教室
	Delete(id string) error
	// Get 取教室；不存在时返回 (nil, nil)
	Get(id string) (*classroom.Classroom, error)
}
