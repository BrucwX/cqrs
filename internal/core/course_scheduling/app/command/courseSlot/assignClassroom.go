package courseSlot

import (
	"context"
	"fmt"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
)

// AssignClassroomInput 给具体课表项安排教室命令
type AssignClassroomInput struct {
	SlotIDs     []string
	ClassroomID string
}

// AssignClassroom 给具体课表项安排教室
//
// 顺序：先判、过了才写。
//
// 规则 = 教室装得下每个目标槽位所属课程的人数 且 与教室现有排课不撞时间。
// 两条判定都交给领域服务：容量走 ClassroomCapacity，时间走 ScheduleConflict。
func (h *Handler) AssignClassroom(ctx context.Context, cmd AssignClassroomInput) (err error) {
	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { err = h.tx.End(ctx, err) }()

	if err := h.CheckClassroom(ctx, &cmd); err != nil {
		return err
	}

	return h.SlotCmd.AssignClassroom(ctx, cmd.SlotIDs, cmd.ClassroomID)
}

func (h *Handler) CheckClassroom(ctx context.Context, cmd *AssignClassroomInput) (err error) {

	// 1) 容量：逐个目标槽位按其所属课程核对教室能不能装下
	tooSmall, err := h.capacity.Check(ctx, cmd.ClassroomID, cmd.SlotIDs)
	if err != nil {
		return err
	}
	if tooSmall {
		return fmt.Errorf("%w: classroom %s", courseSlot.ErrCourseSlotConflict, cmd.ClassroomID)
	}

	// 2) 时间：与教室现有排课重叠 -> 冲突
	//
	// 目标槽位本身已经排了这个教室的话同样会判成冲突 —— 那说明重复安排。
	conflict, err := h.conflict.CheckClassroom(ctx, cmd.ClassroomID, cmd.SlotIDs)
	if err != nil {
		return err
	}
	if conflict {
		return fmt.Errorf("%w: classroom %s", courseSlot.ErrCourseSlotConflict, cmd.ClassroomID)
	}

	return nil
}
