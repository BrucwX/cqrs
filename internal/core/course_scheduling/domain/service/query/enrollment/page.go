package enrollment

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
)

// PageEnrollments 查询输入：课程注册分页列表
type PageEnrollments struct {
	Page     int
	PageSize int
}

// PageEnrollments 课程注册分页查询
func (h *Handler) PageEnrollments(ctx context.Context, q PageEnrollments) ([]*enrollment.CourseEnrollment, error) {
	return h.Query.Page(ctx, q.Page, q.PageSize)
}
