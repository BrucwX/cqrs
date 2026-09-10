package absence

import "cqrs/internal/core/course_scheduling/domain/repo/query"

// Handler 缺勤记录查询处理器
type Handler struct {
	Query query.AbsenceRecordQuery
}
