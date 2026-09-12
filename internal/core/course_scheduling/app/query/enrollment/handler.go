package enrollment

import "cqrs/internal/core/course_scheduling/domain/repo/query"

// Handler 课程注册查询处理器
type Handler struct {
	Query query.CourseEnrollmentQuery
}

// NewHandler 创建课程注册查询处理器
func NewHandler(q query.CourseEnrollmentQuery) *Handler {
	return &Handler{Query: q}
}
