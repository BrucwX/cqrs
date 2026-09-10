package query

import "cqrs/internal/support/commerce/domain/aggregate/enrollment"

// EnrollmentQuery 注册查询接口
type EnrollmentQuery interface {
	// GetByID 根据 ID 获取注册
	GetByID(id string) (*enrollment.Enrollment, error)
	// List 获取注册列表
	List() ([]*enrollment.Enrollment, error)
	// ListByStudentID 根据学员 ID 获取注册列表
	ListByStudentID(studentID string) ([]*enrollment.Enrollment, error)
	// ListByCourseID 根据课程 ID 获取注册列表
	ListByCourseID(courseID string) ([]*enrollment.Enrollment, error)
}
