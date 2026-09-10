package course

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// AvailableForClassroom 查询输入：与教室现有排课不冲突的课程
type AvailableForClassroom struct {
	ClassroomID string
}

// AvailableForClassroomHandler 查询处理器
type AvailableForClassroomHandler struct {
	Query query.CourseQuery
}

func (h *AvailableForClassroomHandler) Execute(ctx context.Context, q AvailableForClassroom) ([]*course.Course, error) {
	return h.Query.AvailableForClassroom(ctx, q.ClassroomID)
}
