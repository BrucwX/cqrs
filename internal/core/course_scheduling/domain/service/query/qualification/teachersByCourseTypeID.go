package qualification

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
)

// TeachersByCourseTypeID 查询输入：根据课程类型 ID 获取有资质的讲师
type TeachersByCourseTypeID struct {
	CourseTypeID string
}

// TeachersByCourseTypeID 根据课程类型 ID 获取有资质的讲师
func (h *Handler) TeachersByCourseTypeID(ctx context.Context, q TeachersByCourseTypeID) ([]*teacher.Teacher, error) {
	return h.Query.TeachersByCourseTypeID(ctx, q.CourseTypeID)
}
