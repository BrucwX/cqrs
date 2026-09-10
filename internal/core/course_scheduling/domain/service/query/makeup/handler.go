package makeup

import "cqrs/internal/core/course_scheduling/domain/repo/query"

// Handler 补课申请查询处理器
type Handler struct {
	Query query.StudentMakeupQuery
}
