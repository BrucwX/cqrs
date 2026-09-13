package implement

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/test_memory/query"
	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

// absenceRecordQuery 是 repoquery.AbsenceRecordQuery 的内存实现。
type absenceRecordQuery struct {
	data *query.Data
}

// NewAbsenceRecordQuery 创建内存版缺勤记录查询。
func NewAbsenceRecordQuery(d *query.Data) repoquery.AbsenceRecordQuery {
	return &absenceRecordQuery{data: d}
}

// Page 分页查询缺勤记录列表。
func (q *absenceRecordQuery) Page(_ context.Context, page, pageSize int) ([]*absence.AbsenceRecord, error) {
	return paginate(q.data.Absences(), page, pageSize), nil
}

// ListByStudentID 根据学员 ID 获取缺勤记录列表。
func (q *absenceRecordQuery) ListByStudentID(_ context.Context, studentID int64) ([]*absence.AbsenceRecord, error) {
	out := make([]*absence.AbsenceRecord, 0)
	for _, item := range q.data.Absences() {
		if item.StudentID() == studentID {
			out = append(out, item)
		}
	}
	return out, nil
}

// ListByCourseID 根据课程 ID 获取缺勤记录列表。
func (q *absenceRecordQuery) ListByCourseID(_ context.Context, courseID string) ([]*absence.AbsenceRecord, error) {
	out := make([]*absence.AbsenceRecord, 0)
	for _, item := range q.data.Absences() {
		if item.CourseID() == courseID {
			out = append(out, item)
		}
	}
	return out, nil
}
