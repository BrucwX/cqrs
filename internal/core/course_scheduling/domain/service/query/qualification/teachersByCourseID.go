package qualification

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// TeachersByCourseID 查询输入：根据课程 ID 获取有资质的讲师
type TeachersByCourseID struct {
	CourseID string
}

// TeachersByCourseIDHandler 查询处理器
type TeachersByCourseIDHandler struct {
	Query query.QualificationQuery
}

func (h *TeachersByCourseIDHandler) Execute(ctx context.Context, q TeachersByCourseID) ([]*teacher.Teacher, error) {
	return h.Query.TeachersByCourseID(ctx, q.CourseID)
}
