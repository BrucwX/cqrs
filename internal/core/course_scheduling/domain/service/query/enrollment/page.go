package enrollment

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// PageEnrollments 查询输入：课程注册分页列表
type PageEnrollments struct {
	Page     int
	PageSize int
}

// PageEnrollmentsHandler 查询处理器
type PageEnrollmentsHandler struct {
	Query query.CourseEnrollmentQuery
}

func (h *PageEnrollmentsHandler) Execute(ctx context.Context, q PageEnrollments) ([]*enrollment.CourseEnrollment, error) {
	return h.Query.Page(ctx, q.Page, q.PageSize)
}
