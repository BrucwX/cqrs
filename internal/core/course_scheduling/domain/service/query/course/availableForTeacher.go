package course

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// AvailableForTeacher 查询输入：与讲师现有排课不冲突的课程
type AvailableForTeacher struct {
	TeacherID int64
}

// AvailableForTeacherHandler 查询处理器
type AvailableForTeacherHandler struct {
	Query query.CourseQuery
}

func (h *AvailableForTeacherHandler) Execute(ctx context.Context, q AvailableForTeacher) ([]*course.Course, error) {
	return h.Query.AvailableForTeacher(ctx, q.TeacherID)
}
