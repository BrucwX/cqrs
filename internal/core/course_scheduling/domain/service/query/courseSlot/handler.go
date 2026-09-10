package courseSlot

import "cqrs/internal/core/course_scheduling/domain/repo/query"

// Handler 课表槽位查询处理器
type Handler struct {
	Query query.CourseSlotQuery
}

// NewHandler 创建课表槽位查询处理器
func NewHandler(q query.CourseSlotQuery) *Handler {
	return &Handler{Query: q}
}
