package student

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// PageStudents 查询输入：学员分页列表
type PageStudents struct {
	Page     int
	PageSize int
}

// PageStudentsHandler 查询处理器
type PageStudentsHandler struct {
	Query query.StudentQuery
}

func (h *PageStudentsHandler) Execute(ctx context.Context, q PageStudents) ([]*student.Student, error) {
	return h.Query.Page(ctx, q.Page, q.PageSize)
}
