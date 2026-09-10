package teacher

import "cqrs/internal/core/course_scheduling/domain/repo/query"

// Handler 讲师查询处理器
type Handler struct {
	Query query.TeacherQuery
}
