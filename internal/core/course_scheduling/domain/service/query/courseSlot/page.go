package courseSlot

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// PageCourseSlots 查询输入：课表槽位分页列表
type PageCourseSlots struct {
	Page     int
	PageSize int
}

// PageCourseSlotsHandler 查询处理器
type PageCourseSlotsHandler struct {
	Query query.CourseSlotQuery
}

func (h *PageCourseSlotsHandler) Execute(ctx context.Context, q PageCourseSlots) ([]*courseSlot.CourseSlot, error) {
	return h.Query.Page(ctx, q.Page, q.PageSize)
}
