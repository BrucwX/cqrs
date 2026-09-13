package query

import (
	"context"
	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

var _ repoquery.ClassroomQuery = (*QueryImpl)(nil)

// PageClassrooms 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) PageClassrooms(ctx context.Context, page, pageSize int) ([]*classroom.Classroom, error) {
	return q.Data.PageClassrooms(ctx, page, pageSize)
}
