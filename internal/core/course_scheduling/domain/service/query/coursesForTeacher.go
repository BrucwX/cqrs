package query

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// CoursesForTeacher 查询输入：与讲师现有排课不冲突的课程
type CoursesForTeacher struct {
	TeacherID int64
}

// CoursesForTeacherHandler 查询处理器
type CoursesForTeacherHandler struct {
	Query query.CourseQuery
}

func (h *CoursesForTeacherHandler) Execute(ctx context.Context, q CoursesForTeacher) ([]*course.Course, error) {
	return h.Query.AvailableForTeacher(ctx, q.TeacherID)
}
