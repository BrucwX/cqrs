package course

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
)

// PageCourses 查询输入：课程分页列表
type PageCourses struct {
	Page     int
	PageSize int
}

// PageCourses 课程分页查询
func (h *Handler) PageCourses(ctx context.Context, q PageCourses) ([]*course.Course, error) {
	return h.Query.Page(ctx, q.Page, q.PageSize)
}
