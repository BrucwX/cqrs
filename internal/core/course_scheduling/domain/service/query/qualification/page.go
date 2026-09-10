package qualification

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// PageQualifications 查询输入：授课资质分页列表
type PageQualifications struct {
	Page     int
	PageSize int
}

// PageQualificationsHandler 查询处理器
type PageQualificationsHandler struct {
	Query query.QualificationQuery
}

func (h *PageQualificationsHandler) Execute(ctx context.Context, q PageQualifications) ([]*qualification.Qualification, error) {
	return h.Query.Page(ctx, q.Page, q.PageSize)
}
