package enrollment

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// CoursesByStudentID 查询输入：根据学员 ID 获取已注册的课程
type CoursesByStudentID struct {
	StudentID int64
}

// CoursesByStudentIDHandler 查询处理器
type CoursesByStudentIDHandler struct {
	Query query.CourseEnrollmentQuery
}

func (h *CoursesByStudentIDHandler) Execute(ctx context.Context, q CoursesByStudentID) ([]*course.Course, error) {
	return h.Query.CoursesByStudentID(ctx, q.StudentID)
}
