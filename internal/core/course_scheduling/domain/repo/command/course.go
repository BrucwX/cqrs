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
	// GetCourse 取课程聚合本身；取不到报 ErrCourseNotFound
	//
	// 与 Get 的差别：Get 找不到时是 (nil, nil)（给「查到了没」的调用方），
	// GetCourse 是规则判定要用的，找不到必须报错。
	GetCourse(courseID string) (course.Course, error)
	// GetCourses 取某课程类型下的全部课程
	GetCourses(courseTypeID string) ([]course.Course, error)
}
