package teacher

import "context"

// DeleteTeacher 删除讲师
//
// 讲师不存在时报 ErrTeacherNotFound。
func (h *Handler) DeleteTeacher(ctx context.Context, id int64) (err error) {
	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { err = h.tx.End(ctx, err) }()

	return h.TeacherCmd.Delete(ctx, id)
}
