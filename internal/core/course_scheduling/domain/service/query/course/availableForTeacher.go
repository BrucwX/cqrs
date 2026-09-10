package course

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
)

// AvailableForTeacher 查询输入：与讲师现有排课不冲突的课程
type AvailableForTeacher struct {
	TeacherID int64
}

// AvailableForTeacher 与讲师现有排课不冲突的课程
func (h *Handler) AvailableForTeacher(ctx context.Context, q AvailableForTeacher) ([]*course.Course, error) {
	return h.Query.AvailableForTeacher(ctx, q.TeacherID)
}
