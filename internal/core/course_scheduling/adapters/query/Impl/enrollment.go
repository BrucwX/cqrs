package query

import (
	"context"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

var _ repoquery.CourseEnrollmentQuery = (*QueryImpl)(nil)

// PageEnrollments 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) PageEnrollments(ctx context.Context, page, pageSize int) ([]*enrollment.CourseEnrollment, error) {
	return q.Data.PageEnrollments(ctx, page, pageSize)
}

// EnrolledCoursesByStudentID 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) EnrolledCoursesByStudentID(ctx context.Context, studentID int64) ([]*course.Course, error) {
	return q.Data.EnrolledCoursesByStudentID(ctx, studentID)
}

// EnrolledStudentsByCourseID 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) EnrolledStudentsByCourseID(ctx context.Context, courseID string) ([]*student.Student, error) {
	return q.Data.EnrolledStudentsByCourseID(ctx, courseID)
}
