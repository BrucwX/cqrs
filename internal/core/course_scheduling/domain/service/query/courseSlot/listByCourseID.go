package courseSlot

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// ListByCourseID 查询输入：根据课程 ID 获取课表槽位
type ListByCourseID struct {
	CourseID string
}

// ListByCourseIDHandler 查询处理器
type ListByCourseIDHandler struct {
	Query query.CourseSlotQuery
}

func (h *ListByCourseIDHandler) Execute(ctx context.Context, q ListByCourseID) ([]*courseSlot.CourseSlot, error) {
	return h.Query.ListByCourseID(q.CourseID)
}
