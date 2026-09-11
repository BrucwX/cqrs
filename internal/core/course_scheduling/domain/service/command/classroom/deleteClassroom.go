package classroom

import "context"

// DeleteClassroom 删除教室
//
// 教室不存在时报 ErrClassroomNotFound。
func (h *Handler) DeleteClassroom(ctx context.Context, id string) error {
	return h.ClassroomCmd.Delete(id)
}
