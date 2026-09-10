package courseSlotChange

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// PageCourseSlotChanges 查询输入：课表变更分页列表
type PageCourseSlotChanges struct {
	Page     int
	PageSize int
}

// PageCourseSlotChangesHandler 查询处理器
type PageCourseSlotChangesHandler struct {
	Query query.CourseSlotChangeQuery
}

func (h *PageCourseSlotChangesHandler) Execute(ctx context.Context, q PageCourseSlotChanges) ([]*courseSlotChange.CourseSlotChange, error) {
	return h.Query.Page(ctx, q.Page, q.PageSize)
}
