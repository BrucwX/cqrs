package absence

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
)

// ListByCourseID 查询输入：根据课程 ID 获取缺勤记录
type ListByCourseID struct {
	CourseID string
}

// ListByCourseID 根据课程 ID 获取缺勤记录
func (h *Handler) ListByCourseID(ctx context.Context, q ListByCourseID) ([]*absence.AbsenceRecord, error) {
	return h.Query.ListAbsencesByCourseID(ctx, q.CourseID)
}
