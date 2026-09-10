package query

import "cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"

// CourseSlotQuery 课表槽位查询接口
type CourseSlotQuery interface {
	// GetByID 根据 ID 获取课表槽位
	GetByID(id int64) (*courseSlot.CourseSlot, error)
	// List 获取课表槽位列表
	List() ([]*courseSlot.CourseSlot, error)
	// ListByCourseID 根据课程 ID 获取课表槽位列表
	ListByCourseID(courseID string) ([]*courseSlot.CourseSlot, error)
}
