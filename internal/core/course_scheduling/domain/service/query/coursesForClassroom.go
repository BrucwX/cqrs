package query

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// CoursesForClassroom 查询输入：与教室现有排课不冲突的课程
type CoursesForClassroom struct {
	ClassroomID string
}

// CoursesForClassroomHandler 查询处理器
type CoursesForClassroomHandler struct {
	Query query.CourseQuery
}

func (h *CoursesForClassroomHandler) Execute(ctx context.Context, q CoursesForClassroom) ([]*course.Course, error) {
	return h.Query.AvailableForClassroom(ctx, q.ClassroomID)
}
