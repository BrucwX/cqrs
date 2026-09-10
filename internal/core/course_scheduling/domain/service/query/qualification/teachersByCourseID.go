package qualification

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
)

// TeachersByCourseID 查询输入：根据课程 ID 获取有资质的讲师
type TeachersByCourseID struct {
	CourseID string
}

// TeachersByCourseID 根据课程 ID 获取有资质的讲师
func (h *Handler) TeachersByCourseID(ctx context.Context, q TeachersByCourseID) ([]*teacher.Teacher, error) {
	return h.Query.TeachersByCourseID(ctx, q.CourseID)
}
