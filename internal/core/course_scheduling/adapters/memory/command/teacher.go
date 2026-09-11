package command

import (
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/memory"
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

// Save 保存讲师（新增或更新）
func (c *TeacherCommand) Save(t *teacher.Teacher) error {
	if t == nil {
		return repo.ErrTeacherRequired
	}
	c.data.SaveTeacher(t)
	return nil
}

// Delete 删除讲师
func (c *TeacherCommand) Delete(id int64) error {
	if !c.data.DeleteTeacher(id) {
		return fmt.Errorf("%w: %d", repo.ErrTeacherNotFound, id)
	}
	return nil
}
