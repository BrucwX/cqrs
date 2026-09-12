package implement

import (
	"context"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/command/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// StudentCommand 学员命令实现（内存版）。
type StudentCommand struct {
	data *memory.Data
}

// 编译期断言：内存实现必须满足命令接口。
var _ repo.StudentCommand = (*StudentCommand)(nil)

// NewStudentCommand 创建学员命令实现
//
// 返回接口类型，这样实现不完整时会在编译期暴露。
func NewStudentCommand(d *memory.Data) repo.StudentCommand {
	return &StudentCommand{data: d}
}

// Create 新增学员
func (c *StudentCommand) Create(ctx context.Context, s *student.Student) error {
	if s == nil {
		return student.ErrStudentRequired
	}
	c.data.SaveStudent(s)
	return nil
}

// Update 按 ID 取出学员交给 updateFn 改，改完写回
func (c *StudentCommand) Update(ctx context.Context, id int64, updateFn func(ctx context.Context, s *student.Student) (*student.Student, error)) error {
	current, err := c.Get(ctx, id)
	if err != nil {
		return err
	}
	if current == nil {
		return fmt.Errorf("%w: %d", student.ErrStudentNotFound, id)
	}

	updated, err := updateFn(ctx, current)
	if err != nil {
		return err
	}
	if updated == nil {
		return student.ErrStudentRequired
	}

	c.data.SaveStudent(updated)
	return nil
}

// Delete 删除学员
func (c *StudentCommand) Delete(ctx context.Context, id int64) error {
	if !c.data.DeleteStudent(id) {
		return fmt.Errorf("%w: %d", student.ErrStudentNotFound, id)
	}
	return nil
}

// Get 取学员；不存在时返回 (nil, nil)。
func (c *StudentCommand) Get(ctx context.Context, id int64) (*student.Student, error) {
	item, ok := c.data.StudentByID(id)
	if !ok {
		return nil, nil
	}
	return item, nil
}

// MustGet 取学员聚合本身；不存在时报 ErrStudentNotFound。
func (c *StudentCommand) MustGet(ctx context.Context, id int64) (student.Student, error) {
	item, ok := c.data.StudentByID(id)
	if !ok {
		return student.Student{}, fmt.Errorf("%w: %d", student.ErrStudentNotFound, id)
	}
	return *item, nil
}
