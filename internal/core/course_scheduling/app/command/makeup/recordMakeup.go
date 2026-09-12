package makeup

import (
	"context"
	"fmt"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
	"cqrs/internal/core/course_scheduling/domain/repo/command"
)

// RecordMakeup 记录学生补课命令
type RecordMakeup struct {
	StudentID      int64
	CourseID       string
	OriginalSlotID string
	OriginalDate   time.Time
	TargetSlotID   string
	TargetDate     time.Time
	MakeupHours    int
}

// RecordMakeup 记录学生补课
//
// 预约即生效（直接落在 StatusBooked），没有审批环节；
// 之后由 CompleteAttendance / Cancel 推进状态。
// 顺序：先判、过了才写。判定交给领域服务 ClassroomCapacity.CheckMakeup，
// 写回走 MakeupCmd.Save。
func (h *Handler) RecordMakeup(ctx context.Context, cmd RecordMakeup) (record *makeup.StudentMakeup, err error) {
	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { err = h.tx.End(ctx, err) }()

	record, err = makeup.NewStudentMakeup(
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

	// 目标那节课的教室得装得下这门课的人 + 已经约在这节课上的其他补课学员 + 这位
	tooSmall, err := h.capacity.CheckMakeup(ctx, cmd.TargetSlotID, cmd.TargetDate)
	if err != nil {
		return nil, err
	}
	if tooSmall {
		return nil, fmt.Errorf("%w: target slot %s", command.ErrMakeupConflict, cmd.TargetSlotID)
	}

	if err := h.MakeupCmd.Save(ctx, record); err != nil {
		return nil, err
	}

	return record, nil
}
