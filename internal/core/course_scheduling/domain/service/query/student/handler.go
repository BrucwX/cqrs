package student

import "cqrs/internal/core/course_scheduling/domain/repo/query"

// Handler 学员查询处理器
type Handler struct {
	Query query.StudentQuery
}
