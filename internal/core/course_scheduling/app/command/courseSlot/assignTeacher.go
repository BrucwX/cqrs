package courseSlot

import (
	"context"
	"fmt"

	"cqrs/internal/core/course_scheduling/domain/repo/command"
)

// AssignTeacherInput 给具体课表项排老师命令
type AssignTeacherInput struct {
	SlotIDs   []string
	TeacherID int64
}

// AssignTeacher 给具体课表项排老师
//
// 顺序：先判、过了才写。判定交给两个领域服务 —— 时间冲突走 ScheduleConflict，
// 资质走 QualificationCheck。
func (h *Handler) AssignTeacher(ctx context.Context, cmd AssignTeacherInput) (err error) {
	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { err = h.tx.End(ctx, err) }()

	// 1) 资质：逐个目标槽位看讲师有没有资格教它所属的那门课
	for _, slotID := range cmd.SlotIDs {
		notQualified, err := h.qualify.CheckTeacherQualifiedForSlot(ctx, cmd.TeacherID, slotID)
		if err != nil {
			return err
		}
		if notQualified {
			return fmt.Errorf("%w: teacher %d 没有该课程类型的资质", command.ErrCourseSlotConflict, cmd.TeacherID)
		}
	}

	// 2) 时间：与讲师现有排课重叠 -> 冲突
	//
	// 目标槽位本身如果已经排了这位讲师，同样会判成冲突 —— 那说明重复安排了，
	// 应该报出来让调用方处理，而不是悄悄放行。
	conflict, err := h.conflict.CheckTeacher(ctx, cmd.TeacherID, cmd.SlotIDs)
	if err != nil {
		return err
	}
	if conflict {
		return fmt.Errorf("%w: teacher %d", command.ErrCourseSlotConflict, cmd.TeacherID)
	}

	return h.SlotCmd.AssignTeacher(ctx, cmd.SlotIDs, cmd.TeacherID)
}
