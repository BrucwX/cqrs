package implement

import (
	"context"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/test_memory/command"
	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// ClassroomCommand 教室命令实现（内存版）。
type ClassroomCommand struct {
	data *command.Data
}

// 编译期断言：内存实现必须满足命令接口。
var _ repo.ClassroomCommand = (*ClassroomCommand)(nil)

// NewClassroomCommand 创建教室命令实现
//
// 返回接口类型，这样实现不完整时会在编译期暴露。
func NewClassroomCommand(d *command.Data) repo.ClassroomCommand {
	return &ClassroomCommand{data: d}
}

// Create 新增教室
func (c *ClassroomCommand) Create(ctx context.Context, cl *classroom.Classroom) error {
	if cl == nil {
		return classroom.ErrClassroomRequired
	}
	c.data.SaveClassroom(cl)
	return nil
}

// Update 按 ID 取出教室交给 updateFn 改，改完写回
//
// 教室不存在时由 MustGet 报 ErrClassroomNotFound。
func (c *ClassroomCommand) Update(ctx context.Context, id string, updateFn func(ctx context.Context, cl *classroom.Classroom) (*classroom.Classroom, error)) error {
	current, err := c.MustGet(ctx, id)
	if err != nil {
		return err
	}

	updated, err := updateFn(ctx, &current)
	if err != nil {
		return err
	}
	if updated == nil {
		return classroom.ErrClassroomRequired
	}

	c.data.SaveClassroom(updated)
	return nil
}

// Delete 删除教室
func (c *ClassroomCommand) Delete(ctx context.Context, id string) error {
	if !c.data.DeleteClassroom(id) {
		return fmt.Errorf("%w: %s", classroom.ErrClassroomNotFound, id)
	}
	return nil
}

// MustGet 取教室聚合本身；不存在时报 ErrClassroomNotFound。
func (c *ClassroomCommand) MustGet(ctx context.Context, id string) (classroom.Classroom, error) {
	item, ok := c.data.ClassroomByID(id)
	if !ok {
		return classroom.Classroom{}, fmt.Errorf("%w: %s", classroom.ErrClassroomNotFound, id)
	}
	return *item, nil
}
