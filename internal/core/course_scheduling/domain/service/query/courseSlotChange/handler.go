package courseSlotChange

import "cqrs/internal/core/course_scheduling/domain/repo/query"

// Handler 课表变更查询处理器
type Handler struct {
	Query query.CourseSlotChangeQuery
}

// NewHandler 创建课表变更查询处理器
func NewHandler(q query.CourseSlotChangeQuery) *Handler {
	return &Handler{Query: q}
}
