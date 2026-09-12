package absence

import "cqrs/internal/core/course_scheduling/domain/repo/query"

// Handler 缺勤记录查询处理器
type Handler struct {
	Query query.AbsenceRecordQuery
}

// NewHandler 创建缺勤记录查询处理器
func NewHandler(q query.AbsenceRecordQuery) *Handler {
	return &Handler{Query: q}
}
