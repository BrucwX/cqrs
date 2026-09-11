package courseSlot

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
)

// AssignTeacherInput 给具体课表项排老师命令
type AssignTeacherInput struct {
	SlotIDs   []string
	TeacherID int64
}

// AssignTeacher 给具体课表项排老师
//
// 仓库（命令适配器）会把「本次要排的槽位 + 该讲师聚合」传进来，判定走 h.checkTeacher。
func (h *Handler) AssignTeacher(ctx context.Context, cmd AssignTeacherInput) error {
	return h.SlotCmd.AssignTeacher(ctx, cmd.SlotIDs, cmd.TeacherID, h.checkTeacher)
}

// checkTeacher 判断把目标槽位排给该讲师是否存在冲突。
//
// 返回 true 表示有冲突（拒绝本次排课），false 表示可以排。
// 规则 = 讲师持有每个目标槽位所属课程类型的资质 且 与讲师现有排课不撞时间。
//
// 仓库只把「本次要排的槽位 + 该讲师聚合」传进来，资质/排期按需自己取。
func (h *Handler) checkTeacher(ctx context.Context, slots []courseSlot.CourseSlot, t teacher.Teacher) (bool, error) {
	// 1) 资质：逐个目标槽位按其所属课程类型核对讲师资质
	qualifications, err := h.assignTea.GetQualifications(t.ID())
	if err != nil {
		return false, err
	}

	for _, slot := range slots {
		courseType, err := h.assignTea.GetCourseType(slot.CourseID())
		if err != nil {
			return false, err
		}

		noQualification, err := checkTeacherQualification(courseType, qualifications)
		if err != nil {
			return false, err
		}
		if noQualification {
			return true, nil
		}
	}

	// 2) 时间：与讲师现有排课重叠 -> 冲突
	//
	// 目标槽位本身如果已经排了这位讲师，同样会判成冲突 —— 那说明重复安排了，
	// 应该报出来让调用方处理，而不是悄悄放行。
	teacherSlots, err := h.assignTea.GetTeacherSlots(t.ID())
	if err != nil {
		return false, err
	}

	return teacherSlots.ConflictsWith(courseSlot.CourseSlots(slots)), nil
}

// checkTeacherQualification 检查老师是否有能力上这门课
// 返回 true 表示老师没有资质（有冲突），false 表示老师有资质（没有冲突）
func checkTeacherQualification(cty courseType.CourseType, t_q []qualification.Qualification) (bool, error) {
	// 检查老师 ID是否匹配
	for _, q := range t_q {
		if q.CourseTypeID() == cty.ID() {
			return false, nil // 老师有资质，没有冲突
		}
	}

	return true, nil // 老师没有资质，有冲突
}
