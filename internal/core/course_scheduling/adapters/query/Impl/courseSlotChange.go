package query

import (
	"context"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

var _ repoquery.CourseSlotChangeQuery = (*QueryImpl)(nil)

// PageCourseSlotChanges 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) PageCourseSlotChanges(ctx context.Context, page, pageSize int) ([]*courseSlotChange.CourseSlotChange, error) {
	return q.Data.PageCourseSlotChanges(ctx, page, pageSize)
}

// ListCourseSlotChangesByCourseID 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) ListCourseSlotChangesByCourseID(ctx context.Context, courseID string) ([]*courseSlotChange.CourseSlotChange, error) {
	return q.Data.ListCourseSlotChangesByCourseID(ctx, courseID)
}
