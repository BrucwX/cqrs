package command

import "cqrs/internal/core/course_scheduling/domain/aggregate/course"

// CourseCommand 课程命令接口
type CourseCommand interface {
	// Save 保存课程（新增或更新）
	Save(c *course.Course) error
	// Delete 删除课程
	Delete(id string) error
}
