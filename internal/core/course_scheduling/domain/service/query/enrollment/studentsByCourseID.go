package enrollment

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// StudentsByCourseID 查询输入：根据课程 ID 获取已注册的学员
type StudentsByCourseID struct {
	CourseID string
}

// StudentsByCourseIDHandler 查询处理器
type StudentsByCourseIDHandler struct {
	Query query.CourseEnrollmentQuery
}

func (h *StudentsByCourseIDHandler) Execute(ctx context.Context, q StudentsByCourseID) ([]*student.Student, error) {
	return h.Query.StudentsByCourseID(ctx, q.CourseID)
}
