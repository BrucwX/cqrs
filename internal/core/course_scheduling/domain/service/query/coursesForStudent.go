package query

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// CoursesForStudent 查询输入：与学员当前选课不冲突的课程
type CoursesForStudent struct {
	StudentID int64
}

// CoursesForStudentHandler 查询处理器
type CoursesForStudentHandler struct {
	Query query.CourseQuery
}

func (h *CoursesForStudentHandler) Execute(ctx context.Context, q CoursesForStudent) ([]*course.Course, error) {
	return h.Query.AvailableForStudent(ctx, q.StudentID)
}
