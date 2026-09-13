package query

import (
	"context"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

var _ repoquery.QualificationQuery = (*QueryImpl)(nil)

// PageQualifications 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) PageQualifications(ctx context.Context, page, pageSize int) ([]*qualification.Qualification, error) {
	return q.Data.PageQualifications(ctx, page, pageSize)
}

// CourseTypesByTeacherID 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) CourseTypesByTeacherID(ctx context.Context, teacherID int64) ([]*courseType.CourseType, error) {
	return q.Data.CourseTypesByTeacherID(ctx, teacherID)
}

// TeachersByCourseTypeID 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) TeachersByCourseTypeID(ctx context.Context, courseTypeID string) ([]*teacher.Teacher, error) {
	return q.Data.TeachersByCourseTypeID(ctx, courseTypeID)
}

// IsTeacherQualifiedForCourse 转发给 Data，等 Redis 接入后在这里组合两者。
func (q *QueryImpl) IsTeacherQualifiedForCourse(ctx context.Context, teacherID int64, courseID string) (bool, error) {
	return q.Data.IsTeacherQualifiedForCourse(ctx, teacherID, courseID)
}
