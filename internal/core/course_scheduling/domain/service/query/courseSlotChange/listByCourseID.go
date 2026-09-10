package courseSlotChange

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// ListByCourseID 查询输入：根据课程 ID 获取课表变更
type ListByCourseID struct {
	CourseID string
}

// ListByCourseIDHandler 查询处理器
type ListByCourseIDHandler struct {
	Query query.CourseSlotChangeQuery
}

func (h *ListByCourseIDHandler) Execute(ctx context.Context, q ListByCourseID) ([]*courseSlotChange.CourseSlotChange, error) {
	return h.Query.ListByCourseID(q.CourseID)
}
