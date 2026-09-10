package course

import "cqrs/internal/core/course_scheduling/domain/repo/query"

// Handler 课程查询处理器
type Handler struct {
	Query query.CourseQuery
}

// NewHandler 创建课程查询处理器
func NewHandler(q query.CourseQuery) *Handler {
	return &Handler{Query: q}
}
