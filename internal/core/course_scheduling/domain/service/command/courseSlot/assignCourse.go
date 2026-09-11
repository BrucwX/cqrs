package courseSlot

import (
	"context"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
)

// AssignCourse 给课程安排课表命令（先不安排老师和教室）
type AssignCourse struct {
	CourseID  string
	Weekday   time.Weekday
	TimeRange courseSlot.DayTimeRange
}

// AssignCourse 给课程安排课表
func (h *Handler) AssignCourse(ctx context.Context, cmd AssignCourse) error {
	// 保存
	return h.SlotCmd.AssignCourse(ctx, []string{cmd.CourseID}, cmd.CourseID, checkScheduleConflict)
}
