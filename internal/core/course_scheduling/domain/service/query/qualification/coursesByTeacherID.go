package qualification

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
)

// CoursesByTeacherID 查询输入：根据讲师 ID 获取有资质的课程
type CoursesByTeacherID struct {
	TeacherID int64
}

// CoursesByTeacherID 根据讲师 ID 获取有资质的课程
func (h *Handler) CoursesByTeacherID(ctx context.Context, q CoursesByTeacherID) ([]*course.Course, error) {
	return h.Query.CoursesByTeacherID(ctx, q.TeacherID)
}
