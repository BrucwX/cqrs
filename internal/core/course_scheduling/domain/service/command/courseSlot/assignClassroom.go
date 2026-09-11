package courseSlot

import (
	"context"
)

// AssignClassroom 给具体课表项安排教室命令
type AssignClassroom struct {
	SlotIDs     []int64
	ClassroomID string
}

// AssignClassroom 给具体课表项安排教室
func (h *Handler) AssignClassroom(ctx context.Context, cmd AssignClassroom) error {
	return h.SlotCmd.AssignClassroom(ctx, cmd.SlotIDs, cmd.ClassroomID, checkScheduleConflict)
}
