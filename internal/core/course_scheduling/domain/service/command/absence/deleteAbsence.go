package absence

import "context"

// DeleteAbsence 删除缺勤记录
//
// 记录不存在时报 ErrAbsenceNotFound。
func (h *Handler) DeleteAbsence(ctx context.Context, id int64) error {
	return h.AbsenceCmd.Delete(id)
}
