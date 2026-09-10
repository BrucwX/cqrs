package courseSlotChange

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
)

// PageCourseSlotChanges 查询输入：课表变更分页列表
type PageCourseSlotChanges struct {
	Page     int
	PageSize int
}

// PageCourseSlotChanges 课表变更分页查询
func (h *Handler) PageCourseSlotChanges(ctx context.Context, q PageCourseSlotChanges) ([]*courseSlotChange.CourseSlotChange, error) {
	return h.Query.Page(ctx, q.Page, q.PageSize)
}
