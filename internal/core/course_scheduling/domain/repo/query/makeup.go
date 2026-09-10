package query

import "cqrs/internal/core/course_scheduling/domain/aggregate/makeup"

// StudentMakeupQuery 补课申请查询接口
type StudentMakeupQuery interface {
	// GetByID 根据 ID 获取补课申请
	GetByID(id int64) (*makeup.StudentMakeup, error)
	// List 获取补课申请列表
	List() ([]*makeup.StudentMakeup, error)
	// ListByStudentID 根据学员 ID 获取补课申请列表
	ListByStudentID(studentID int64) ([]*makeup.StudentMakeup, error)
	// ListByCourseID 根据课程 ID 获取补课申请列表
	ListByCourseID(courseID string) ([]*makeup.StudentMakeup, error)
}
