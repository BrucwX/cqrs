package student

import "context"

// DeleteStudent 删除学员
//
// 学员不存在时报 ErrStudentNotFound。
func (h *Handler) DeleteStudent(ctx context.Context, id int64) error {
	return h.StudentCmd.Delete(id)
}
