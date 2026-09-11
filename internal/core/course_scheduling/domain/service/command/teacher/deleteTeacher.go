package teacher

import "context"

// DeleteTeacher 删除讲师
//
// 讲师不存在时报 ErrTeacherNotFound。
func (h *Handler) DeleteTeacher(ctx context.Context, id int64) error {
	return h.TeacherCmd.Delete(id)
}
