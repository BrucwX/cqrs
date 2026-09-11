package courseSlotChange

import (
	"context"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
)

// ChangeCourseSlot 临时换课命令
type ChangeCourseSlot struct {
	CourseID    string
	ApplicantID int64
	ChangeType  courseSlotChange.ChangeType
	Original    courseSlotChange.OriginalPlan
	Target      courseSlotChange.TargetPlan
	Reason      string
}

// ChangeCourseSlot 临时换课（调课 / 代课 / 换教室）
//
// 登记即生效，没有审批环节；参数校验由聚合根构造函数负责
// （课程/申请人/事由必填，目标时间不能早于现在，目标起止时间与讲师、教室必填）。
// 仓库（命令适配器）会把「本次换课」传进来，冲突判定走 h.checkSlotChange。
func (h *Handler) ChangeCourseSlot(ctx context.Context, cmd ChangeCourseSlot) (*courseSlotChange.CourseSlotChange, error) {
	change, err := courseSlotChange.NewCourseSlotChange(
		cmd.CourseID,
		cmd.ApplicantID,
		cmd.ChangeType,
		cmd.Original,
		cmd.Target,
		cmd.Reason,
		time.Now(),
	)
	if err != nil {
		return nil, err
	}

	if err := h.ChangeCmd.Change(ctx, change, h.checkSlotChange); err != nil {
		return nil, err
	}

	return change, nil
}

// checkSlotChange 判断这次换课是否存在冲突。
//
// 返回 true 表示有冲突（拒绝本次换课），false 表示可以换。
// 规则 = 目标讲师在目标时段没别的课 且 目标教室在目标时段没被占用
// 且 没有别的换课已经把该讲师/该教室占在同一时间段。
//
// 仓库只把「本次换课」传进来，讲师/教室排期、其他换课记录都按需自己取。
func (h *Handler) checkSlotChange(ctx context.Context, csc *courseSlotChange.CourseSlotChange) (bool, error) {
	target := csc.TargetPlan()

	teacherSlots, err := h.slotChange.GetTeacherSlots(target.TeacherID())
	if err != nil {
		return false, err
	}

	classroomSlots, err := h.slotChange.GetClassroomSlots(target.ClassroomID())
	if err != nil {
		return false, err
	}

	// 1) 周排期：目标时段落在目标讲师 / 目标教室已有的课表上。
	//    本门课自己的排期不算 —— 那就是被换掉的那节课本身。
	span, err := daySpan(target)
	if err != nil {
		return false, err
	}

	weekday := target.TargetStartAt().Weekday()
	if otherCourseSlots(teacherSlots, csc.CourseID()).OverlapsOn(weekday, span) {
		return true, nil // 讲师那个时间已经在给别人上课
	}
	if otherCourseSlots(classroomSlots, csc.CourseID()).OverlapsOn(weekday, span) {
		return true, nil // 教室那个时间已经被别的课占用
	}

	// 2) 换课记录：别的临时换课已经把该讲师 / 该教室占在同一时间段
	others, err := h.slotChange.GetOtherSlotChanges(csc.ID())
	if err != nil {
		return false, err
	}

	for _, other := range others {
		ot := other.TargetPlan()
		if !overlaps(
			ot.TargetStartAt(), ot.TargetEndAt(),
			target.TargetStartAt(), target.TargetEndAt(),
		) {
			continue
		}
		if ot.TeacherID() == target.TeacherID() || ot.ClassroomID() == target.ClassroomID() {
			return true, nil
		}
	}

	return false, nil
}

// otherCourseSlots 剔除属于该课程自己的排期 —— 那就是被换掉的那节课本身，不算占用。
func otherCourseSlots(slots courseSlot.CourseSlots, courseID string) courseSlot.CourseSlots {
	out := make(courseSlot.CourseSlots, 0, len(slots))
	for _, cs := range slots {
		if cs.CourseID() != courseID {
			out = append(out, cs)
		}
	}
	return out
}

// daySpan 把目标时间段折算成「当天几点到几点」。
//
// TargetPlan 保证起止落在同一天内，所以正常不会失败；
// 万一从存储里恢复出跨天记录，这里直接报错而不是静默跳过比对。
func daySpan(target courseSlotChange.TargetPlan) (courseSlot.DayTimeRange, error) {
	start, end := target.TargetStartAt(), target.TargetEndAt()

	from, err := courseSlot.NewDayTime(start.Hour(), start.Minute())
	if err != nil {
		return courseSlot.DayTimeRange{}, err
	}
	to, err := courseSlot.NewDayTime(end.Hour(), end.Minute())
	if err != nil {
		return courseSlot.DayTimeRange{}, err
	}

	return courseSlot.NewDayTimeRange(from, to)
}

// overlaps 两个具体时间段是否有交集。
func overlaps(aStart, aEnd, bStart, bEnd time.Time) bool {
	return aStart.Before(bEnd) && bStart.Before(aEnd)
}
