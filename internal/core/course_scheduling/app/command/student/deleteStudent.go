package student

import "context"

// DeleteStudent 删除学员
//
// 学员不存在时报 ErrStudentNotFound。
func (h *Handler) DeleteStudent(ctx context.Context, id int64) (err error) {
	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { err = h.tx.End(ctx, err) }()

	return h.StudentCmd.Delete(id)
}
