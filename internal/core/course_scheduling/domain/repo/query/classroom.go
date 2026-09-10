package query

import "cqrs/internal/core/course_scheduling/domain/aggregate/classroom"

// ClassroomQuery 教室查询接口
type ClassroomQuery interface {
	// GetByID 根据 ID 获取教室
	GetByID(id string) (*classroom.Classroom, error)
	// List 获取教室列表
	List() ([]*classroom.Classroom, error)
}
