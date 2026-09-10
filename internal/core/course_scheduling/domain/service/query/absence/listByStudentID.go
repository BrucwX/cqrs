package absence

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
)

// ListByStudentID 查询输入：根据学员 ID 获取缺勤记录
type ListByStudentID struct {
	StudentID int64
}

// ListByStudentID 根据学员 ID 获取缺勤记录
func (h *Handler) ListByStudentID(ctx context.Context, q ListByStudentID) ([]*absence.AbsenceRecord, error) {
	return h.Query.ListByStudentID(q.StudentID)
}
