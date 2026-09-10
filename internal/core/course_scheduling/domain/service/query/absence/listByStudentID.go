package absence

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// ListByStudentID 查询输入：根据学员 ID 获取缺勤记录
type ListByStudentID struct {
	StudentID int64
}

// ListByStudentIDHandler 查询处理器
type ListByStudentIDHandler struct {
	Query query.AbsenceRecordQuery
}

func (h *ListByStudentIDHandler) Execute(ctx context.Context, q ListByStudentID) ([]*absence.AbsenceRecord, error) {
	return h.Query.ListByStudentID(q.StudentID)
}
