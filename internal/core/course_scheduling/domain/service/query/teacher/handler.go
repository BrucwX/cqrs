package teacher

import "cqrs/internal/core/course_scheduling/domain/repo/query"

// Handler 讲师查询处理器
type Handler struct {
	Query query.TeacherQuery
}

// NewHandler 创建讲师查询处理器
func NewHandler(q query.TeacherQuery) *Handler {
	return &Handler{Query: q}
}
