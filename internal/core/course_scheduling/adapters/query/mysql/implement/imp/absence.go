package imp

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/query/mysql"
	"cqrs/internal/core/course_scheduling/adapters/query/mysql/implement/help"
	"cqrs/internal/core/course_scheduling/adapters/query/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

// absenceRecordQuery 是 repoquery.AbsenceRecordQuery 的 MySQL 实现。
type absenceRecordQuery struct {
	data *mysql.Data
}

// 编译期断言：实现必须满足接口。
var _ repoquery.AbsenceRecordQuery = (*absenceRecordQuery)(nil)

// NewAbsenceRecordQuery 创建 MySQL 版缺勤记录查询。
func NewAbsenceRecordQuery(d *mysql.Data) repoquery.AbsenceRecordQuery {
	return &absenceRecordQuery{data: d}
}

// Page 分页查询缺勤记录列表，按 ID 升序。
func (q *absenceRecordQuery) Page(ctx context.Context, page, pageSize int) ([]*absence.AbsenceRecord, error) {
	limit, offset := help.LimitOffset(page, pageSize)
	return help.QueryAll(ctx, q.data.Conn(ctx), `
SELECT `+help.AbsenceColumns+`
  FROM absence_record
 ORDER BY id
 LIMIT ? OFFSET ?`, []any{limit, offset}, help.ScanAbsence, model.AbsenceToDO)
}

// ListByStudentID 根据学员 ID 获取缺勤记录列表。
func (q *absenceRecordQuery) ListByStudentID(ctx context.Context, studentID int64) ([]*absence.AbsenceRecord, error) {
	return help.QueryAll(ctx, q.data.Conn(ctx), `
SELECT `+help.AbsenceColumns+`
  FROM absence_record
 WHERE student_id = ?
 ORDER BY id`, []any{studentID}, help.ScanAbsence, model.AbsenceToDO)
}

// ListByCourseID 根据课程 ID 获取缺勤记录列表。
func (q *absenceRecordQuery) ListByCourseID(ctx context.Context, courseID string) ([]*absence.AbsenceRecord, error) {
	return help.QueryAll(ctx, q.data.Conn(ctx), `
SELECT `+help.AbsenceColumns+`
  FROM absence_record
 WHERE course_id = ?
 ORDER BY id`, []any{courseID}, help.ScanAbsence, model.AbsenceToDO)
}
