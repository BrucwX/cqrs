package qualification

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
)

// CourseTypesByTeacherID 查询输入：根据讲师 ID 获取有资质的课程类型
//
// 资质绑定的是课程类型（而不是具体课程），所以返回的是课程类型列表。
type CourseTypesByTeacherID struct {
	TeacherID int64
}

// CourseTypesByTeacherID 根据讲师 ID 获取有资质的课程类型
func (h *Handler) CourseTypesByTeacherID(ctx context.Context, q CourseTypesByTeacherID) ([]*courseType.CourseType, error) {
	return h.Query.CourseTypesByTeacherID(ctx, q.TeacherID)
}
