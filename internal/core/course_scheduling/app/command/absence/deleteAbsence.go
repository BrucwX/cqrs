package absence

import "context"

// DeleteAbsence 删除缺勤记录
//
// 记录不存在时报 ErrAbsenceNotFound。
func (h *Handler) DeleteAbsence(ctx context.Context, id int64) (err error) {
	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { err = h.tx.End(ctx, err) }()

	return h.AbsenceCmd.Delete(id)
}
