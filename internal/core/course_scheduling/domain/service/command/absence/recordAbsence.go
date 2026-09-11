package absence

import (
	"context"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
)

// RecordAbsence 记录学生缺课命令
type RecordAbsence struct {
	StudentID    int64
	CourseID     string
	CourseSlotID int64
	ScheduleDate time.Time
	MissedHours  int
	AbsenceType  absence.AbsenceType
	Reason       string
}

// RecordAbsence 记录学生缺课
//
// 缺课就是一条事实记录，没有审批也没有「是否已补卡」这类状态；
// 事假/公假/旷课的区分交给 AbsenceType。
func (h *Handler) RecordAbsence(ctx context.Context, cmd RecordAbsence) (*absence.AbsenceRecord, error) {
	record, err := absence.NewAbsenceRecord(
		cmd.StudentID,
		cmd.CourseID,
		cmd.CourseSlotID,
		cmd.ScheduleDate,
		cmd.MissedHours,
		cmd.AbsenceType,
		cmd.Reason,
	)
	if err != nil {
		return nil, err
	}

	if err := h.AbsenceCmd.Save(record); err != nil {
		return nil, err
	}

	return record, nil
}
