package classroom

import "context"

// DeleteClassroom 删除教室
//
// 教室不存在时报 ErrClassroomNotFound。
func (h *Handler) DeleteClassroom(ctx context.Context, id string) (err error) {
	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { err = h.tx.End(ctx, err) }()

	return h.ClassroomCmd.DeleteClassroom(ctx, id)
}
