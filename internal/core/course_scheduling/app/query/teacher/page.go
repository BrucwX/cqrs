package teacher

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
)

// PageTeachers 查询输入：讲师分页列表
type PageTeachers struct {
	Page     int
	PageSize int
}

// PageTeachers 讲师分页查询
func (h *Handler) PageTeachers(ctx context.Context, q PageTeachers) ([]*teacher.Teacher, error) {
	return h.Query.PageTeachers(ctx, q.Page, q.PageSize)
}
