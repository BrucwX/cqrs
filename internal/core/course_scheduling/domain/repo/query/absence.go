package query

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
)

// AbsenceRecordQuery 缺勤记录查询接口
type AbsenceRecordQuery interface {
	// PageAbsences 分页查询缺勤记录列表
	PageAbsences(ctx context.Context, page int, pageSize int) ([]*absence.AbsenceRecord, error)
	// ListByStudentID 根据学员 ID 获取缺勤记录列表
	ListByStudentID(ctx context.Context, studentID int64) ([]*absence.AbsenceRecord, error)
	// ListAbsencesByCourseID 根据课程 ID 获取缺勤记录列表
	ListAbsencesByCourseID(ctx context.Context, courseID string) ([]*absence.AbsenceRecord, error)
}
