package courseSlot

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
)

// PageCourseSlots 查询输入：课表槽位分页列表
type PageCourseSlots struct {
	Page     int
	PageSize int
}

// PageCourseSlots 课表槽位分页查询
func (h *Handler) PageCourseSlots(ctx context.Context, q PageCourseSlots) ([]*courseSlot.CourseSlot, error) {
	return h.Query.PageCourseSlots(ctx, q.Page, q.PageSize)
}
