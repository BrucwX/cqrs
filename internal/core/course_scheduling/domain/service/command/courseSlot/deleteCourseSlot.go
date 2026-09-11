package courseSlot

import "context"

// DeleteCourseSlot 删除课表槽位
//
// 槽位不存在时报 ErrCourseSlotNotFound。
func (h *Handler) DeleteCourseSlot(ctx context.Context, id string) error {
	return h.SlotCmd.Delete(id)
}
