package query

import "cqrs/internal/core/course_scheduling/domain/aggregate/qualification"

// QualificationQuery 授课资质查询接口
type QualificationQuery interface {
	// GetByID 根据 ID 获取授课资质
	GetByID(id int64) (*qualification.Qualification, error)
	// List 获取授课资质列表
	List() ([]*qualification.Qualification, error)
	// ListByTeacherID 根据讲师 ID 获取授课资质列表
	ListByTeacherID(teacherID int64) ([]*qualification.Qualification, error)
	// ListByCourseID 根据课程 ID 获取授课资质列表
	ListByCourseID(courseID string) ([]*qualification.Qualification, error)
}
