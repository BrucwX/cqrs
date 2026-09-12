package command

import (
	"errors"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
)

var (
	// ErrCourseTypeNotFound 课程或它归属的课程类型不存在。
	ErrCourseTypeNotFound = errors.New("course type not found")
)

// CourseTypeCommand 课程类型命令接口
//
// 目前只有「按课程取它归属的类型」这一个读方法 —— 排讲师时要核对资质、
// 授资质时要知道课程属于哪个类型。课程类型本身还没有增删改的用例。
type CourseTypeCommand interface {
	// GetCourseType 取某门课程归属的课程类型
	GetCourseType(courseID string) (courseType.CourseType, error)
}
