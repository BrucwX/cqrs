package query

import (
	"context"
	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

var _ repoquery.StudentQuery = (*QueryImpl)(nil)

// PageStudents 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) PageStudents(ctx context.Context, page, pageSize int) ([]*student.Student, error) {
	return q.Data.PageStudents(ctx, page, pageSize)
}
