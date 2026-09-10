package qualification

import "cqrs/internal/core/course_scheduling/domain/repo/query"

// Handler 授课资质查询处理器
type Handler struct {
	Query query.QualificationQuery
}
