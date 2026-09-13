package courseSlotChange

import "context"

// DeleteCourseSlotChange 删除课表变更单
//
// 变更单不存在时报 ErrSlotChangeNotFound。
func (h *Handler) DeleteCourseSlotChange(ctx context.Context, id int64) (err error) {
	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { err = h.tx.End(ctx, err) }()

	return h.ChangeCmd.DeleteCourseSlotChange(ctx, id)
}
