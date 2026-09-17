package query

import (
	"context"
	"courseScheduling/internal/domain/aggregate/classroom"
	repoquery "courseScheduling/internal/domain/repo/query"
)

var _ repoquery.ClassroomQuery = (*QueryImpl)(nil)

// PageClassrooms 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) PageClassrooms(ctx context.Context, page, pageSize int) ([]*classroom.Classroom, error) {
	return q.Data.PageClassrooms(ctx, page, pageSize)
}
