package query

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
)

// AbsenceRecordQuery 缺勤记录查询接口
type AbsenceRecordQuery interface {
	// Page 分页查询缺勤记录列表
	Page(ctx context.Context, page int, pageSize int) ([]*absence.AbsenceRecord, error)
	// ListByStudentID 根据学员 ID 获取缺勤记录列表
	ListByStudentID(studentID int64) ([]*absence.AbsenceRecord, error)
	// ListByCourseID 根据课程 ID 获取缺勤记录列表
	ListByCourseID(courseID string) ([]*absence.AbsenceRecord, error)
}
