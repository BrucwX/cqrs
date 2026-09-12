package makeup

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
)

// PageMakeups 查询输入：补课申请分页列表
type PageMakeups struct {
	Page     int
	PageSize int
}

// PageMakeups 补课申请分页查询
func (h *Handler) PageMakeups(ctx context.Context, q PageMakeups) ([]*makeup.StudentMakeup, error) {
	return h.Query.Page(ctx, q.Page, q.PageSize)
}
