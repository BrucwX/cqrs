package enrollment

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
)

// StudentsByCourseID 查询输入：根据课程 ID 获取已注册的学员
type StudentsByCourseID struct {
	CourseID string
}

// StudentsByCourseID 根据课程 ID 获取已注册的学员
func (h *Handler) StudentsByCourseID(ctx context.Context, q StudentsByCourseID) ([]*student.Student, error) {
	return h.Query.StudentsByCourseID(ctx, q.CourseID)
}
