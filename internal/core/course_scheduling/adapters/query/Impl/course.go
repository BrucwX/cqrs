package query

import (
	"context"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

var _ repoquery.CourseQuery = (*QueryImpl)(nil)

// PageCourses 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) PageCourses(ctx context.Context, page, pageSize int) ([]*course.Course, error) {
	return q.Data.PageCourses(ctx, page, pageSize)
}

// AvailableForStudent 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) AvailableForStudent(ctx context.Context, studentID int64) ([]*course.Course, error) {
	return q.Data.AvailableForStudent(ctx, studentID)
}

// AvailableForTeacher 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) AvailableForTeacher(ctx context.Context, teacherID int64) ([]*course.Course, error) {
	return q.Data.AvailableForTeacher(ctx, teacherID)
}

// AvailableForClassroom 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) AvailableForClassroom(ctx context.Context, classroomID string) ([]*course.Course, error) {
	return q.Data.AvailableForClassroom(ctx, classroomID)
}
