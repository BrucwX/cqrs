package enrollment

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
)

// CoursesByStudentID 查询输入：根据学员 ID 获取已注册的课程
type CoursesByStudentID struct {
	StudentID int64
}

// CoursesByStudentID 根据学员 ID 获取已注册的课程
func (h *Handler) CoursesByStudentID(ctx context.Context, q CoursesByStudentID) ([]*course.Course, error) {
	return h.Query.CoursesByStudentID(ctx, q.StudentID)
}
