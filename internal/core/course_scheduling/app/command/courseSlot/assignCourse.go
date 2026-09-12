package courseSlot

import (
	"context"
	"fmt"

	"cqrs/internal/core/course_scheduling/domain/repo/command"
)

// AssignCourseInput 给指定课表槽位设置课程命令
type AssignCourseInput struct {
	SlotIDs  []string
	CourseID string
}

// AssignCourse 给指定课表槽位设置课程
//
// 顺序：先判、过了才写。判定交给领域服务 ScheduleConflict
// （目标课程在不在、时间撞不撞都由它一次判完）。
func (h *Handler) AssignCourse(ctx context.Context, cmd AssignCourseInput) (err error) {
	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { err = h.tx.End(ctx, err) }()

	conflict, err := h.conflict.CheckCourse(ctx, cmd.CourseID, cmd.SlotIDs)
	if err != nil {
		return err
	}
	if conflict {
		return fmt.Errorf("%w: course %s", command.ErrCourseSlotConflict, cmd.CourseID)
	}

	return h.SlotCmd.AssignCourse(ctx, cmd.SlotIDs, cmd.CourseID)
}
