package absence

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// ListByCourseID 查询输入：根据课程 ID 获取缺勤记录
type ListByCourseID struct {
	CourseID string
}

// ListByCourseIDHandler 查询处理器
type ListByCourseIDHandler struct {
	Query query.AbsenceRecordQuery
}

func (h *ListByCourseIDHandler) Execute(ctx context.Context, q ListByCourseID) ([]*absence.AbsenceRecord, error) {
	return h.Query.ListByCourseID(q.CourseID)
}
