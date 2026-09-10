package query

import "cqrs/internal/core/course_scheduling/domain/aggregate/course"

// CourseQuery 课程查询接口
type CourseQuery interface {
	// GetByID 根据 ID 获取课程
	GetByID(id string) (*course.Course, error)
	// List 获取课程列表
	List() ([]*course.Course, error)
}
