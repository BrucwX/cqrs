package command

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
)

// CourseSlotCommand 课表槽位命令实现
type CourseSlotCommand struct {
	data *memory.Data
}

// NewCourseSlotCommand 创建课表槽位命令实现
func NewCourseSlotCommand(d *memory.Data) *CourseSlotCommand {
	return &CourseSlotCommand{data: d}
}

// Save 保存课表槽位
func (c *CourseSlotCommand) Save(cs *courseSlot.CourseSlot) error {
	// TODO: 实现保存逻辑
	return nil
}

// Delete 删除课表槽位
func (c *CourseSlotCommand) Delete(id int64) error {
	// TODO: 实现删除逻辑
	return nil
}

// Update 更新课表槽位
func (c *CourseSlotCommand) Update(ctx context.Context, slotIDs []int64, teacherID int64) error {
	// TODO: 实现更新逻辑
	return nil
}

// AssignTeacher 给指定课表槽位们配置老师
func (c *CourseSlotCommand) AssignTeacher(ctx context.Context, slotIDs []int64, teacherID int64, checkConflictFn func(ctx context.Context, cs []courseSlot.CourseSlot, t_cs []courseSlot.CourseSlot) (bool, error)) error {
	// TODO: 实现配置老师逻辑
	return nil
}

// AssignCourse 给指定课表槽位们配置课程
func (c *CourseSlotCommand) AssignCourse(ctx context.Context, slotIDs []string, courseId string, checkConflictFn func(ctx context.Context, cs []courseSlot.CourseSlot, c_cs []courseSlot.CourseSlot) (bool, error)) error {
	// TODO: 实现配置课程逻辑
	return nil
}

// AssignClassroom 给指定课表槽位们配置教室
func (c *CourseSlotCommand) AssignClassroom(ctx context.Context, slotIDs []int64, classroomID string, checkConflictFn func(ctx context.Context, cs []courseSlot.CourseSlot, cl_cs []courseSlot.CourseSlot) (bool, error)) error {
	// TODO: 实现配置教室逻辑
	return nil
}
