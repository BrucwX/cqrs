package command

import (
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/memory"
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

// Save 保存学员（新增或更新）
func (c *StudentCommand) Save(s *student.Student) error {
	if s == nil {
		return repo.ErrStudentRequired
	}
	c.data.SaveStudent(s)
	return nil
}

// Delete 删除学员
func (c *StudentCommand) Delete(id int64) error {
	if !c.data.DeleteStudent(id) {
		return fmt.Errorf("%w: %d", repo.ErrStudentNotFound, id)
	}
	return nil
}
