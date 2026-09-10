package query

import "cqrs/internal/core/course_scheduling/domain/aggregate/teacher"

// TeacherQuery 讲师查询接口
type TeacherQuery interface {
	// GetByID 根据 ID 获取讲师
	GetByID(id int64) (*teacher.Teacher, error)
	// List 获取讲师列表
	List() ([]*teacher.Teacher, error)
}
