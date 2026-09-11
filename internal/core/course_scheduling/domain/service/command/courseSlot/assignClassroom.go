package courseSlot

import (
	"context"

	command "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// AssignClassroom 给具体课表项安排教室命令
type AssignClassroom struct {
	SlotIDs     []string
	ClassroomID string
}

// AssignClassroom 给具体课表项安排教室
//
// 仓库（命令适配器）会把 checkClassroom 需要的数据装进 ClassroomAssignContext
// 一并传进来，所以这里直接把 checkClassroom 当回调交给它。
func (h *Handler) AssignClassroom(ctx context.Context, cmd AssignClassroom) error {
	return h.SlotCmd.AssignClassroom(ctx, cmd.SlotIDs, cmd.ClassroomID, checkClassroom)
}

// checkClassroom 判断把目标槽位排给该教室是否存在冲突。
//
// 返回 true 表示有冲突（拒绝本次排课），false 表示可以排。
// 规则 = 教室装得下课程人数 且 与教室现有排课不撞时间。
func checkClassroom(ctx context.Context, ac command.ClassroomAssignContext) (bool, error) {
	// 1) 容量：直接用教室自带的方法（它还会判教室状态与已分配座位）
	if err := ac.Classroom.CanAccommodate(ac.Course.Capacity().Max()); err != nil {
		return true, nil // 装不下 / 不可用 -> 冲突
	}

	// 2) 时间：与教室现有排课重叠 -> 冲突
	return ac.TargetSlots.ConflictsWith(ac.ClassroomSlots), nil
}
