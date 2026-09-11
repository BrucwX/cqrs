package command

import (
	"cqrs/internal/core/course_scheduling/adapters/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
)

// EnrollmentCommand 课程注册命令实现
type EnrollmentCommand struct {
	data *memory.Data
}

// NewEnrollmentCommand 创建课程注册命令实现
func NewEnrollmentCommand(d *memory.Data) *EnrollmentCommand {
	return &EnrollmentCommand{data: d}
}

// Save 保存课程注册
func (c *EnrollmentCommand) Save(e *enrollment.CourseEnrollment) error {
	// TODO: 实现保存逻辑
	return nil
}

// Delete 删除课程注册
func (c *EnrollmentCommand) Delete(id int64) error {
	// TODO: 实现删除逻辑
	return nil
}
