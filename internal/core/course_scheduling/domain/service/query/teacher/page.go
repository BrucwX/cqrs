package teacher

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
	"cqrs/internal/core/course_scheduling/domain/repo/query"
)

// PageTeachers 查询输入：讲师分页列表
type PageTeachers struct {
	Page     int
	PageSize int
}

// PageTeachersHandler 查询处理器
type PageTeachersHandler struct {
	Query query.TeacherQuery
}

func (h *PageTeachersHandler) Execute(ctx context.Context, q PageTeachers) ([]*teacher.Teacher, error) {
	return h.Query.Page(ctx, q.Page, q.PageSize)
}
