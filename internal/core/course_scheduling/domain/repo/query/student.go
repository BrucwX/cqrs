package query

import "cqrs/internal/core/course_scheduling/domain/aggregate/student"

// StudentQuery 学员查询接口
type StudentQuery interface {
	// GetByID 根据 ID 获取学员
	GetByID(id int64) (*student.Student, error)
	// List 获取学员列表
	List() ([]*student.Student, error)
}
