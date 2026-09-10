package courseSlot

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
)

// ListByCourseID 查询输入：根据课程 ID 获取课表槽位
type ListByCourseID struct {
	CourseID string
}

// ListByCourseID 根据课程 ID 获取课表槽位
func (h *Handler) ListByCourseID(ctx context.Context, q ListByCourseID) ([]*courseSlot.CourseSlot, error) {
	return h.Query.ListByCourseID(q.CourseID)
}
