package command

import "cqrs/internal/core/course_scheduling/domain/aggregate/classroom"

// ClassroomCommand 教室命令接口
type ClassroomCommand interface {
	// Save 保存教室（新增或更新）
	Save(c *classroom.Classroom) error
	// Delete 删除教室
	Delete(id string) error
}
