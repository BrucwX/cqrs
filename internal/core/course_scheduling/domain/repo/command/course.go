package command

import (
	"errors"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
)

var (
	// ErrCourseRequired 传入的课程为空。
	ErrCourseRequired = errors.New("course is required")
	// ErrCourseNotFound 指定的课程不存在。
	ErrCourseNotFound = errors.New("course not found")
)

// CourseCommand 课程命令接口
type CourseCommand interface {
	// Save 保存课程（新增或更新）
	Save(c *course.Course) error
	// Delete 删除课程
	Delete(id string) error
}
