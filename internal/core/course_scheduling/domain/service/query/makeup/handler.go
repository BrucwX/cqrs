package makeup

import "cqrs/internal/core/course_scheduling/domain/repo/query"

// Handler 补课申请查询处理器
type Handler struct {
	Query query.StudentMakeupQuery
}

// NewHandler 创建补课申请查询处理器
func NewHandler(q query.StudentMakeupQuery) *Handler {
	return &Handler{Query: q}
}
