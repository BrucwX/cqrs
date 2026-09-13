package student

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
)

// PageStudents 查询输入：学员分页列表
type PageStudents struct {
	Page     int
	PageSize int
}

// PageStudents 学员分页查询
func (h *Handler) PageStudents(ctx context.Context, q PageStudents) ([]*student.Student, error) {
	return h.Query.PageStudents(ctx, q.Page, q.PageSize)
}
