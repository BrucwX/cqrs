package command

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
)

// CourseCommand 课程命令接口
type CourseCommand interface {
	// CreateCourse 新增课程
	CreateCourse(ctx context.Context, c *course.Course) error
	// UpdateCourse 更新课程：按 ID 取出已有课程交给 updateFn 改，改完写回
	//
	// 课程不存在时报 ErrCourseNotFound。
	UpdateCourse(
		ctx context.Context,
		id string,
		updateFn func(ctx context.Context, crs *course.Course) (*course.Course, error),
	) error
	// DeleteCourse 删除课程
	DeleteCourse(ctx context.Context, id string) error
	// GetCourse 取课程；不存在时返回 (nil, nil)
	GetCourse(ctx context.Context, id string) (*course.Course, error)
	// MustGetCourse 取课程聚合本身；取不到报 ErrCourseNotFound
	//
	// 与 Get 的差别：Get 找不到时是 (nil, nil)（给「查到了没」的调用方），
	// MustGetCourse 是规则判定要用的，找不到必须报错。
	MustGetCourse(ctx context.Context, id string) (course.Course, error)
	// GetCourses 取某课程类型下的全部课程
	GetCourses(ctx context.Context, courseTypeID string) ([]course.Course, error)
}
