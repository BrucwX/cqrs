package makeup

import (
	"context"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
)

// RecordMakeup 记录学生补课命令
type RecordMakeup struct {
	StudentID      int64
	CourseID       string
	OriginalSlotID int64
	OriginalDate   time.Time
	TargetSlotID   int64
	TargetDate     time.Time
	MakeupHours    int
}

// RecordMakeup 记录学生补课
//
// 预约即生效（直接落在 StatusBooked），没有审批环节；
// 之后由 CompleteAttendance / Cancel 推进状态。
func (h *Handler) RecordMakeup(ctx context.Context, cmd RecordMakeup) (*makeup.StudentMakeup, error) {
	record, err := makeup.NewStudentMakeup(
		cmd.StudentID,
		cmd.CourseID,
		cmd.OriginalSlotID,
		cmd.OriginalDate,
		cmd.TargetSlotID,
		cmd.TargetDate,
		cmd.MakeupHours,
		time.Now(),
	)
	if err != nil {
		return nil, err
	}

	if err := h.MakeupCmd.Save(record); err != nil {
		return nil, err
	}

	return record, nil
}
