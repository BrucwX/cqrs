package qualification

import "context"

// IsTeacherQualifiedForCourse 查询输入：讲师是否有教某门课的资质
type IsTeacherQualifiedForCourse struct {
	TeacherID int64
	CourseID  string
}

// IsTeacherQualifiedForCourse 讲师是否有教某门课的资质
//
// 链路：课程 → 所属课程类型 → 该讲师是否有该类型的资质。
func (h *Handler) IsTeacherQualifiedForCourse(ctx context.Context, q IsTeacherQualifiedForCourse) (bool, error) {
	return h.Query.IsTeacherQualifiedForCourse(ctx, q.TeacherID, q.CourseID)
}
