package scheduleConflict

import (
	"context"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
)

// CheckSlotChange 判断这次临时换课是否存在冲突。
//
// 规则 = 目标讲师在目标时段没别的课 且 目标教室在目标时段没被占用
// 且 没有别的换课已经把该讲师 / 该教室占在同一时间段。
//
// 前两条看的是周排期（CourseSlot），第三条看的是还没落成排期的换课单
// （CourseSlotChange）。换课单本身填得对不对（课程 / 申请人 / 事由必填、目标
// 时间不能早于现在…）归 courseSlotChange 聚合根的构造函数管，不在这里。
func (s *Service) CheckSlotChange(ctx context.Context, csc *courseSlotChange.CourseSlotChange) (bool, error) {
	target := csc.TargetPlan()

	teacherSlots, err := s.SlotCmd.GetTeacherSlots(target.TeacherID())
	if err != nil {
		return false, err
	}

	classroomSlots, err := s.SlotCmd.GetClassroomSlots(target.ClassroomID())
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
	others, err := s.changes.GetOtherSlotChanges(csc.ID())
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
