package mysql

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/query/data/mysql/help"
	"cqrs/internal/core/course_scheduling/adapters/query/data/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
)

// Page 分页查询缺勤记录列表，按 ID 升序。
func (q *Data) PageAbsences(ctx context.Context, page, pageSize int) ([]*absence.AbsenceRecord, error) {
	limit, offset := help.LimitOffset(page, pageSize)
	return help.QueryAll(ctx, q.Conn(ctx), `
SELECT `+help.AbsenceColumns+`
  FROM absence_record
 ORDER BY id
 LIMIT ? OFFSET ?`, []any{limit, offset}, help.ScanAbsence, model.AbsenceToDO)
}

// ListByStudentID 根据学员 ID 获取缺勤记录列表。
func (q *Data) ListByStudentID(ctx context.Context, studentID int64) ([]*absence.AbsenceRecord, error) {
	return help.QueryAll(ctx, q.Conn(ctx), `
SELECT `+help.AbsenceColumns+`
  FROM absence_record
 WHERE student_id = ?
 ORDER BY id`, []any{studentID}, help.ScanAbsence, model.AbsenceToDO)
}

// ListByCourseID 根据课程 ID 获取缺勤记录列表。
func (q *Data) ListAbsencesByCourseID(ctx context.Context, courseID string) ([]*absence.AbsenceRecord, error) {
	return help.QueryAll(ctx, q.Conn(ctx), `
SELECT `+help.AbsenceColumns+`
  FROM absence_record
 WHERE course_id = ?
 ORDER BY id`, []any{courseID}, help.ScanAbsence, model.AbsenceToDO)
}
