package implement

import (
	"context"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/memImp4test/command"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// TeacherCommand 讲师命令实现（内存版）。
type TeacherCommand struct {
	data *command.Data
}

// 编译期断言：内存实现必须满足命令接口。
var _ repo.TeacherCommand = (*TeacherCommand)(nil)

// NewTeacherCommand 创建讲师命令实现
//
// 返回接口类型，这样实现不完整时会在编译期暴露。
func NewTeacherCommand(d *command.Data) repo.TeacherCommand {
	return &TeacherCommand{data: d}
}

// CreateTeacher 新增讲师
func (c *TeacherCommand) CreateTeacher(ctx context.Context, t *teacher.Teacher) error {
	if t == nil {
		return teacher.ErrTeacherRequired
	}
	c.data.SaveTeacher(t)
	return nil
}

// UpdateTeacher 按 ID 取出讲师交给 updateFn 改，改完写回
func (c *TeacherCommand) UpdateTeacher(ctx context.Context, id int64, updateFn func(ctx context.Context, t *teacher.Teacher) (*teacher.Teacher, error)) error {
	current, err := c.GetTeacher(ctx, id)
	if err != nil {
		return err
	}
	if current == nil {
		return fmt.Errorf("%w: %d", teacher.ErrTeacherNotFound, id)
	}

	updated, err := updateFn(ctx, current)
	if err != nil {
		return err
	}
	if updated == nil {
		return teacher.ErrTeacherRequired
	}

	c.data.SaveTeacher(updated)
	return nil
}

// DeleteTeacher 删除讲师
func (c *TeacherCommand) DeleteTeacher(ctx context.Context, id int64) error {
	if !c.data.DeleteTeacher(id) {
		return fmt.Errorf("%w: %d", teacher.ErrTeacherNotFound, id)
	}
	return nil
}

// GetTeacher 取讲师；不存在时返回 (nil, nil)。
func (c *TeacherCommand) GetTeacher(ctx context.Context, id int64) (*teacher.Teacher, error) {
	item, ok := c.data.TeacherByID(id)
	if !ok {
		return nil, nil
	}
	return item, nil
}

// MustGetTeacher 取讲师聚合本身；取不到报 ErrTeacherNotFound。
//
// 与 Get 的差别：Get 找不到时是 (nil, nil)（给「查到了没」的调用方），
// MustGetTeacher 是规则判定要用的，找不到必须报错。
func (c *TeacherCommand) MustGetTeacher(ctx context.Context, teacherID int64) (teacher.Teacher, error) {
	item, ok := c.data.TeacherByID(teacherID)
	if !ok {
		return teacher.Teacher{}, fmt.Errorf("%w: %d", teacher.ErrTeacherNotFound, teacherID)
	}
	return *item, nil
}
