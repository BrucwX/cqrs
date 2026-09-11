package courseSlot

import (
	"context"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
)

// ScheduleCourse 给课程安排课表（先不安排老师和教室）
type ScheduleCourse struct {
	CourseID  string
	Weekday   time.Weekday
	TimeRange courseSlot.DayTimeRange
}

func (h *Handler) ScheduleCourse(ctx context.Context, cmd ScheduleCourse) error {
	cs, err := courseSlot.NewCourseSlot(
		cmd.CourseID,
		cmd.Weekday,
		cmd.TimeRange,
		courseSlot.PendingTeacherID,
		courseSlot.PendingClassroomID,
	)
	if err != nil {
		return err
	}
	return h.SlotCmd.Save(cs)
}
