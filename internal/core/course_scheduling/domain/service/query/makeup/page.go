package makeup

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// PageMakeups 查询输入：补课申请分页列表
type PageMakeups struct {
	Page     int
	PageSize int
}

// PageMakeupsHandler 查询处理器
type PageMakeupsHandler struct {
	Query query.StudentMakeupQuery
}

func (h *PageMakeupsHandler) Execute(ctx context.Context, q PageMakeups) ([]*makeup.StudentMakeup, error) {
	return h.Query.Page(ctx, q.Page, q.PageSize)
}
