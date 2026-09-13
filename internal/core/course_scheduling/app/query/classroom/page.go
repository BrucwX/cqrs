package classroom

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
)

// PageClassrooms 查询输入：教室分页列表
type PageClassrooms struct {
	Page     int
	PageSize int
}

// PageClassrooms 教室分页查询
func (h *Handler) PageClassrooms(ctx context.Context, q PageClassrooms) ([]*classroom.Classroom, error) {
	return h.Query.PageClassrooms(ctx, q.Page, q.PageSize)
}
