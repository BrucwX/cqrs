package query

import "cqrs/internal/core/course_scheduling/domain/aggregate/absence"

// AbsenceRecordQuery 缺勤记录查询接口
type AbsenceRecordQuery interface {
	// GetByID 根据 ID 获取缺勤记录
	GetByID(id int64) (*absence.AbsenceRecord, error)
	// List 获取缺勤记录列表
	List() ([]*absence.AbsenceRecord, error)
	// ListByStudentID 根据学员 ID 获取缺勤记录列表
	ListByStudentID(studentID int64) ([]*absence.AbsenceRecord, error)
	// ListByCourseID 根据课程 ID 获取缺勤记录列表
	ListByCourseID(courseID string) ([]*absence.AbsenceRecord, error)
}
