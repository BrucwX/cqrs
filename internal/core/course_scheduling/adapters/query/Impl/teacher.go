package query

import (
	"context"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

var _ repoquery.TeacherQuery = (*QueryImpl)(nil)

// PageTeachers 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) PageTeachers(ctx context.Context, page, pageSize int) ([]*teacher.Teacher, error) {
	return q.Data.PageTeachers(ctx, page, pageSize)
}
