package absence

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// PageAbsences 查询输入：缺勤记录分页列表
type PageAbsences struct {
	Page     int
	PageSize int
}

// PageAbsencesHandler 查询处理器
type PageAbsencesHandler struct {
	Query query.AbsenceRecordQuery
}

func (h *PageAbsencesHandler) Execute(ctx context.Context, q PageAbsences) ([]*absence.AbsenceRecord, error) {
	return h.Query.Page(ctx, q.Page, q.PageSize)
}
