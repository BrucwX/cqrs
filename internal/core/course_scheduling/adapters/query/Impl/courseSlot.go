package query

import (
	"context"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

var _ repoquery.CourseSlotQuery = (*QueryImpl)(nil)

// PageCourseSlots 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) PageCourseSlots(ctx context.Context, page, pageSize int) ([]*courseSlot.CourseSlot, error) {
	return q.Data.PageCourseSlots(ctx, page, pageSize)
}

// ListCourseSlotsByCourseID 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) ListCourseSlotsByCourseID(ctx context.Context, courseID string) ([]*courseSlot.CourseSlot, error) {
	return q.Data.ListCourseSlotsByCourseID(ctx, courseID)
}
