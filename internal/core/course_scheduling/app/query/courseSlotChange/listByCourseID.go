package courseSlotChange

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
)

// ListByCourseID 查询输入：根据课程 ID 获取课表变更
type ListByCourseID struct {
	CourseID string
}

// ListByCourseID 根据课程 ID 获取课表变更
func (h *Handler) ListByCourseID(ctx context.Context, q ListByCourseID) ([]*courseSlotChange.CourseSlotChange, error) {
	return h.Query.ListCourseSlotChangesByCourseID(ctx, q.CourseID)
}
