package absence

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
)

// PageAbsences 查询输入：缺勤记录分页列表
type PageAbsences struct {
	Page     int
	PageSize int
}

// PageAbsences 缺勤记录分页查询
func (h *Handler) PageAbsences(ctx context.Context, q PageAbsences) ([]*absence.AbsenceRecord, error) {
	return h.Query.PageAbsences(ctx, q.Page, q.PageSize)
}
