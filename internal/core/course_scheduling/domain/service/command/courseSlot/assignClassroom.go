package courseSlot

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
)

// AssignClassroomInput 给具体课表项安排教室命令
type AssignClassroomInput struct {
	SlotIDs     []string
	ClassroomID string
}

// AssignClassroom 给具体课表项安排教室
//
// 仓库（命令适配器）会把「本次要排的槽位 + 该教室聚合」传进来，判定走 h.checkClassroom。
func (h *Handler) AssignClassroom(ctx context.Context, cmd AssignClassroomInput) error {
	return h.SlotCmd.AssignClassroom(ctx, cmd.SlotIDs, cmd.ClassroomID, h.checkClassroom)
}

// checkClassroom 判断把目标槽位排给该教室是否存在冲突。
//
// 返回 true 表示有冲突（拒绝本次排课），false 表示可以排。
// 规则 = 教室装得下每个目标槽位所属课程的人数 且 与教室现有排课不撞时间。
//
// 目标槽位本身已经排了这个教室的话同样会判成冲突 —— 那说明重复安排。
//
// 仓库只把「本次要排的槽位 + 该教室聚合」传进来，课程/教室排期按需自己取。
func (h *Handler) checkClassroom(ctx context.Context, slots []courseSlot.CourseSlot, c classroom.Classroom) (bool, error) {
	// 1) 容量：逐个目标槽位按其所属课程核对教室能不能装下
	for _, slot := range slots {
		crs, err := h.assignClass.GetCourse(slot.CourseID())
		if err != nil {
			return false, err
		}

		// 教室自带的方法，还会判教室状态与已分配座位
		if err := c.CanAccommodate(crs.Capacity().Max()); err != nil {
			return true, nil // 装不下 / 不可用 -> 冲突
		}
	}

	// 2) 时间：与教室现有排课重叠 -> 冲突
	classroomSlots, err := h.assignClass.GetClassroomSlots(c.ID())
	if err != nil {
		return false, err
	}

	return classroomSlots.ConflictsWith(courseSlot.CourseSlots(slots)), nil
}
