package implement

import (
	"context"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/command/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// TeacherCommand 讲师命令实现（内存版）。
type TeacherCommand struct {
	data *memory.Data
}

// 编译期断言：内存实现必须满足命令接口。
var _ repo.TeacherCommand = (*TeacherCommand)(nil)

// NewTeacherCommand 创建讲师命令实现
//
// 返回接口类型，这样实现不完整时会在编译期暴露。
func NewTeacherCommand(d *memory.Data) repo.TeacherCommand {
	return &TeacherCommand{data: d}
}

// Create 新增讲师
func (c *TeacherCommand) Create(t *teacher.Teacher) error {
	if t == nil {
		return repo.ErrTeacherRequired
	}
	c.data.SaveTeacher(t)
	return nil
}

// Update 按 ID 取出讲师交给 updateFn 改，改完写回
func (c *TeacherCommand) Update(ctx context.Context, id int64, updateFn func(ctx context.Context, t *teacher.Teacher) (*teacher.Teacher, error)) error {
	current, err := c.Get(id)
	if err != nil {
		return err
	}
	if current == nil {
		return fmt.Errorf("%w: %d", repo.ErrTeacherNotFound, id)
	}

	updated, err := updateFn(ctx, current)
	if err != nil {
		return err
	}
	if updated == nil {
		return repo.ErrTeacherRequired
	}

	c.data.SaveTeacher(updated)
	return nil
}

// Delete 删除讲师
func (c *TeacherCommand) Delete(id int64) error {
	if !c.data.DeleteTeacher(id) {
		return fmt.Errorf("%w: %d", repo.ErrTeacherNotFound, id)
	}
	return nil
}

// Get 取讲师；不存在时返回 (nil, nil)。
func (c *TeacherCommand) Get(id int64) (*teacher.Teacher, error) {
	item, ok := c.data.TeacherByID(id)
	if !ok {
		return nil, nil
	}
	return item, nil
}
