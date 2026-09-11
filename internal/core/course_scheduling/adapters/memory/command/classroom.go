package command

import (
	"context"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// ClassroomCommand 教室命令实现（内存版）。
type ClassroomCommand struct {
	data *memory.Data
}

// 编译期断言：内存实现必须满足命令接口。
var _ repo.ClassroomCommand = (*ClassroomCommand)(nil)

// NewClassroomCommand 创建教室命令实现
//
// 返回接口类型，这样实现不完整时会在编译期暴露。
func NewClassroomCommand(d *memory.Data) repo.ClassroomCommand {
	return &ClassroomCommand{data: d}
}

// Create 新增教室
func (c *ClassroomCommand) Create(cl *classroom.Classroom) error {
	if cl == nil {
		return repo.ErrClassroomRequired
	}
	c.data.SaveClassroom(cl)
	return nil
}

// Update 按 ID 取出教室交给 updateFn 改，改完写回
func (c *ClassroomCommand) Update(ctx context.Context, id string, updateFn func(ctx context.Context, cl *classroom.Classroom) (*classroom.Classroom, error)) error {
	current, err := c.Get(id)
	if err != nil {
		return err
	}
	if current == nil {
		return fmt.Errorf("%w: %s", repo.ErrClassroomNotFound, id)
	}

	updated, err := updateFn(ctx, current)
	if err != nil {
		return err
	}
	if updated == nil {
		return repo.ErrClassroomRequired
	}

	c.data.SaveClassroom(updated)
	return nil
}

// Delete 删除教室
func (c *ClassroomCommand) Delete(id string) error {
	if !c.data.DeleteClassroom(id) {
		return fmt.Errorf("%w: %s", repo.ErrClassroomNotFound, id)
	}
	return nil
}

// Get 取教室；不存在时返回 (nil, nil)，由调用方决定怎么处理。
func (c *ClassroomCommand) Get(id string) (*classroom.Classroom, error) {
	item, ok := c.data.ClassroomByID(id)
	if !ok {
		return nil, nil
	}
	return item, nil
}
