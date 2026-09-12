package makeup

import "context"

// DeleteMakeup 删除补课预约
//
// 预约不存在时报 ErrMakeupNotFound。注意这是硬删除，
// 需要保留取消痕迹请改用 StudentMakeup.Cancel。
func (h *Handler) DeleteMakeup(ctx context.Context, id int64) (err error) {
	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { err = h.tx.End(ctx, err) }()

	return h.MakeupCmd.Delete(id)
}
