package courseSlot

import "context"

// DeleteCourseSlot 删除课表槽位
//
// 槽位不存在时报 ErrCourseSlotNotFound。
func (h *Handler) DeleteCourseSlot(ctx context.Context, id string) (err error) {
	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { err = h.tx.End(ctx, err) }()

	return h.SlotCmd.Delete(id)
}
