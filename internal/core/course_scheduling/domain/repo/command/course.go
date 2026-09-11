package command

import (
	"context"
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
	// Create 新增课程
	Create(c *course.Course) error
	// Update 更新课程：按 ID 取出已有课程交给 updateFn 改，改完写回
	//
	// 课程不存在时报 ErrCourseNotFound。
	Update(
		ctx context.Context,
		id string,
		updateFn func(ctx context.Context, crs *course.Course) (*course.Course, error),
	) error
	// Delete 删除课程
	Delete(id string) error
	// Get 取课程；不存在时返回 (nil, nil)
	Get(id string) (*course.Course, error)
}
