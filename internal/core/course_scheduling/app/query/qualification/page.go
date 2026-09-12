package qualification

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
)

// PageQualifications 查询输入：授课资质分页列表
type PageQualifications struct {
	Page     int
	PageSize int
}

// PageQualifications 授课资质分页查询
func (h *Handler) PageQualifications(ctx context.Context, q PageQualifications) ([]*qualification.Qualification, error) {
	return h.Query.Page(ctx, q.Page, q.PageSize)
}
