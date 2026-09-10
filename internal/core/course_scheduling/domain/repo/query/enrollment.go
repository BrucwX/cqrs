package query

import "cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"

// CourseEnrollmentQuery 课程注册查询接口
type CourseEnrollmentQuery interface {
	// GetByID 根据 ID 获取课程注册
	GetByID(id int64) (*enrollment.CourseEnrollment, error)
	// List 获取课程注册列表
	List() ([]*enrollment.CourseEnrollment, error)
	// ListByStudentID 根据学员 ID 获取课程注册列表
	ListByStudentID(studentID int64) ([]*enrollment.CourseEnrollment, error)
	// ListByCourseID 根据课程 ID 获取课程注册列表
	ListByCourseID(courseID string) ([]*enrollment.CourseEnrollment, error)
}
