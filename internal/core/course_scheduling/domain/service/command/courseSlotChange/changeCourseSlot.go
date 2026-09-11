package courseSlotChange

import (
	"context"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
	command "cqrs/internal/core/course_scheduling/domain/repo/command"
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
// 冲突判定走 checkSlotChange：仓库（命令适配器）会把 SlotChangeContext 装好传进来。
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

	if err := h.ChangeCmd.Change(ctx, change, checkSlotChange); err != nil {
		return nil, err
	}

	return change, nil
}

// checkSlotChange 判断这次换课是否存在冲突。
//
// 返回 true 表示有冲突（拒绝本次换课），false 表示可以换。
// 规则 = 目标讲师在目标时段没别的课 且 目标教室在目标时段没被占用
// 且 没有别的换课已经把该讲师/该教室占在同一时间段。
func checkSlotChange(ctx context.Context, sc command.SlotChangeContext) (bool, error) {
	target := sc.Change.TargetPlan()

	// 1) 周排期：目标时段落在目标讲师 / 目标教室已有的课表上
	//    目标跨天时无法折算成「当天的几点到几点」，跳过这一条，只比换课记录。
	if span, ok := daySpan(target); ok {
		weekday := target.TargetStartAt().Weekday()
		if sc.TeacherSlots.OverlapsOn(weekday, span) {
			return true, nil // 讲师那个时间已经在给别人上课
		}
		if sc.ClassroomSlots.OverlapsOn(weekday, span) {
			return true, nil // 教室那个时间已经被别的课占用
		}
	}

	// 2) 换课记录：别的临时换课已经把该讲师 / 该教室占在同一时间段
	for _, other := range sc.OtherChanges {
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

// daySpan 把目标时间段折算成「当天几点到几点」。
//
// 目标跨天（或时间点非法）时返回 false，调用方只能退回按具体时间比对。
func daySpan(target courseSlotChange.TargetPlan) (courseSlot.DayTimeRange, bool) {
	start, end := target.TargetStartAt(), target.TargetEndAt()
	if start.Year() != end.Year() || start.YearDay() != end.YearDay() {
		return courseSlot.DayTimeRange{}, false
	}

	from, err := courseSlot.NewDayTime(start.Hour(), start.Minute())
	if err != nil {
		return courseSlot.DayTimeRange{}, false
	}
	to, err := courseSlot.NewDayTime(end.Hour(), end.Minute())
	if err != nil {
		return courseSlot.DayTimeRange{}, false
	}

	span, err := courseSlot.NewDayTimeRange(from, to)
	if err != nil {
		return courseSlot.DayTimeRange{}, false
	}
	return span, true
}

// overlaps 两个具体时间段是否有交集。
func overlaps(aStart, aEnd, bStart, bEnd time.Time) bool {
	return aStart.Before(bEnd) && bStart.Before(aEnd)
}
