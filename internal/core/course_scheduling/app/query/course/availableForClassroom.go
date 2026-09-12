package course

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
)

// AvailableForClassroom 查询输入：与教室现有排课不冲突的课程
type AvailableForClassroom struct {
	ClassroomID string
}

// AvailableForClassroom 与教室现有排课不冲突的课程
func (h *Handler) AvailableForClassroom(ctx context.Context, q AvailableForClassroom) ([]*course.Course, error) {
	return h.Query.AvailableForClassroom(ctx, q.ClassroomID)
}
