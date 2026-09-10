package course

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// AvailableForStudent 查询输入：与学员当前选课不冲突的课程
type AvailableForStudent struct {
	StudentID int64
}

// AvailableForStudentHandler 查询处理器
type AvailableForStudentHandler struct {
	Query query.CourseQuery
}

func (h *AvailableForStudentHandler) Execute(ctx context.Context, q AvailableForStudent) ([]*course.Course, error) {
	return h.Query.AvailableForStudent(ctx, q.StudentID)
}
