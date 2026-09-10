package query

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
)

// QualificationQuery 授课资质查询接口
type QualificationQuery interface {
	// Page 分页查询授课资质列表
	Page(ctx context.Context, page int, pageSize int) ([]*qualification.Qualification, error)
	// CoursesByTeacherID 根据讲师 ID 获取有资质的课程列表
	CoursesByTeacherID(ctx context.Context, teacherID int64) ([]*course.Course, error)
	// TeachersByCourseID 根据课程 ID 获取有资质的讲师列表
	TeachersByCourseID(ctx context.Context, courseID string) ([]*teacher.Teacher, error)
}
