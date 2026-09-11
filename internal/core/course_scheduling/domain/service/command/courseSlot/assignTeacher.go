package courseSlot

import (
	"context"

	command "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// AssignTeacher 给具体课表项排老师命令
type AssignTeacher struct {
	SlotIDs   []string
	TeacherID int64
}

// AssignTeacher 给具体课表项排老师
//
// 仓库（命令适配器）会把 checkTeacher 需要的数据装进 TeacherAssignContext
// 一并传进来，所以这里直接把 checkTeacher 当回调交给它。
func (h *Handler) AssignTeacher(ctx context.Context, cmd AssignTeacher) error {
	return h.SlotCmd.AssignTeacher(ctx, cmd.SlotIDs, cmd.TeacherID, checkTeacher)
}

// checkTeacher 判断把目标槽位排给该讲师是否存在冲突。
//
// 返回 true 表示有冲突（拒绝本次排课），false 表示可以排。
// 规则 = 讲师有该课程类型的资质 且 与讲师现有排课不撞时间。
func checkTeacher(ctx context.Context, ac command.TeacherAssignContext) (bool, error) {
	// 1) 资质：讲师没有该课程类型的资质 -> 冲突
	noQualification, err := checkTeacherQualification(ctx, ac.CourseType, ac.Qualifications)
	if err != nil {
		return false, err
	}
	if noQualification {
		return true, nil
	}

	// 2) 时间：与讲师现有排课重叠 -> 冲突
	return ac.TargetSlots.ConflictsWith(ac.TeacherSlots), nil
}
