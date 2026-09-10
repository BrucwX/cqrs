package classroom

import "cqrs/internal/core/course_scheduling/domain/repo/query"

// Handler 教室查询处理器
type Handler struct {
	Query query.ClassroomQuery
}
