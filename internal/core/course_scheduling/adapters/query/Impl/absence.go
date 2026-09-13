package query

import (
	"context"
	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

var _ repoquery.AbsenceRecordQuery = (*QueryImpl)(nil)

// PageAbsences 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) PageAbsences(ctx context.Context, page, pageSize int) ([]*absence.AbsenceRecord, error) {
	return q.Data.PageAbsences(ctx, page, pageSize)
}

// ListByStudentID 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) ListByStudentID(ctx context.Context, studentID int64) ([]*absence.AbsenceRecord, error) {
	return q.Data.ListByStudentID(ctx, studentID)
}

// ListAbsencesByCourseID 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) ListAbsencesByCourseID(ctx context.Context, courseID string) ([]*absence.AbsenceRecord, error) {
	return q.Data.ListAbsencesByCourseID(ctx, courseID)
}
