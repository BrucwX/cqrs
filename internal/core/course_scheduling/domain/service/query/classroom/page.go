package classroom

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// PageClassrooms 查询输入：教室分页列表
type PageClassrooms struct {
	Page     int
	PageSize int
}

// PageClassroomsHandler 查询处理器
type PageClassroomsHandler struct {
	Query query.ClassroomQuery
}

func (h *PageClassroomsHandler) Execute(ctx context.Context, q PageClassrooms) ([]*classroom.Classroom, error) {
	return h.Query.Page(ctx, q.Page, q.PageSize)
}
