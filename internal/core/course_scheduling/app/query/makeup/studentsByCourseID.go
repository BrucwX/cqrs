package makeup

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
)

// StudentsByCourseID 查询输入：根据课程 ID 获取已补课的学员
type StudentsByCourseID struct {
	CourseID string
}

// StudentsByCourseID 根据课程 ID 获取已补课的学员
func (h *Handler) StudentsByCourseID(ctx context.Context, q StudentsByCourseID) ([]*student.Student, error) {
	return h.Query.MakeupStudentsByCourseID(ctx, q.CourseID)
}
