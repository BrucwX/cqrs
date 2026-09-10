package query

import "cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"

// CourseSlotChangeQuery 课表变更查询接口
type CourseSlotChangeQuery interface {
	// GetByID 根据 ID 获取课表变更
	GetByID(id int64) (*courseSlotChange.CourseSlotChange, error)
	// List 获取课表变更列表
	List() ([]*courseSlotChange.CourseSlotChange, error)
	// ListByCourseID 根据课程 ID 获取课表变更列表
	ListByCourseID(courseID string) ([]*courseSlotChange.CourseSlotChange, error)
}
