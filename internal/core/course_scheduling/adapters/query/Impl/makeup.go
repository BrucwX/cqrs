package query

import (
	"context"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

var _ repoquery.StudentMakeupQuery = (*QueryImpl)(nil)

// PageMakeups 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) PageMakeups(ctx context.Context, page, pageSize int) ([]*makeup.StudentMakeup, error) {
	return q.Data.PageMakeups(ctx, page, pageSize)
}

// MakeupCoursesByStudentID 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) MakeupCoursesByStudentID(ctx context.Context, studentID int64) ([]*course.Course, error) {
	return q.Data.MakeupCoursesByStudentID(ctx, studentID)
}

// MakeupStudentsByCourseID 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) MakeupStudentsByCourseID(ctx context.Context, courseID string) ([]*student.Student, error) {
	return q.Data.MakeupStudentsByCourseID(ctx, courseID)
}
