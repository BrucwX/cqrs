package courseSlotChange

import "context"

// DeleteCourseSlotChange 删除课表变更单
//
// 变更单不存在时报 ErrSlotChangeNotFound。
func (h *Handler) DeleteCourseSlotChange(ctx context.Context, id int64) error {
	return h.ChangeCmd.Delete(id)
}
