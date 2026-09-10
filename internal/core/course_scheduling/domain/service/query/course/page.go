package course

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// PageCourses 查询输入：课程分页列表
type PageCourses struct {
	Page     int
	PageSize int
}

// PageCoursesHandler 查询处理器
type PageCoursesHandler struct {
	Query query.CourseQuery
}

func (h *PageCoursesHandler) Execute(ctx context.Context, q PageCourses) ([]*course.Course, error) {
	return h.Query.Page(ctx, q.Page, q.PageSize)
}
