package courseSlot

import (
	"context"
)

// AssignTeacher 给具体课表项排老师命令
type AssignTeacher struct {
	SlotIDs   []int64
	TeacherID int64
}

// AssignTeacher 给具体课表项排老师
func (h *Handler) AssignTeacher(ctx context.Context, cmd AssignTeacher) error {
	return h.SlotCmd.AssignTeacher(ctx, cmd.SlotIDs, cmd.TeacherID)
}
