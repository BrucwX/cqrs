package command

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
)

// CourseTypeCommand 课程类型命令接口
//
// 目前只有「按课程取它归属的类型」这一个读方法 —— 排讲师时要核对资质、
// 授资质时要知道课程属于哪个类型。课程类型本身还没有增删改的用例。
type CourseTypeCommand interface {
	// MustGetCourseType 取课程类型聚合；不存在时报 ErrCourseTypeNotFound
	MustGetCourseType(ctx context.Context, id string) (courseType.CourseType, error)
	// GetCourseType 取某门课程归属的课程类型
	GetCourseType(ctx context.Context, courseID string) (courseType.CourseType, error)
}
