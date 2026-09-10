package qualification

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// CoursesByTeacherID 查询输入：根据讲师 ID 获取有资质的课程
type CoursesByTeacherID struct {
	TeacherID int64
}

// CoursesByTeacherIDHandler 查询处理器
type CoursesByTeacherIDHandler struct {
	Query query.QualificationQuery
}

func (h *CoursesByTeacherIDHandler) Execute(ctx context.Context, q CoursesByTeacherID) ([]*course.Course, error) {
	return h.Query.CoursesByTeacherID(ctx, q.TeacherID)
}
