package course

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
)

// AvailableForStudent 查询输入：与学员当前选课不冲突的课程
type AvailableForStudent struct {
	StudentID int64
}

// AvailableForStudent 与学员当前选课不冲突的课程
func (h *Handler) AvailableForStudent(ctx context.Context, q AvailableForStudent) ([]*course.Course, error) {
	return h.Query.AvailableForStudent(ctx, q.StudentID)
}
