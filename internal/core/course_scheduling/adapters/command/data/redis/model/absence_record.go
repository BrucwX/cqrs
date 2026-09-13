package model

import (
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
)

// AbsenceRecord 是 Redis 快照（对应 absence_record），领域模型是 absence.AbsenceRecord。
//
// 一条事实：没有审批、也没有「已补卡」这类状态，所以聚合里没有任何状态。
type AbsenceRecord struct {
	ID           int64     `json:"id"`             // 缺勤记录 ID
	StudentID    int64     `json:"student_id"`     // 学员/员工 ID
	CourseID     string    `json:"course_id"`      // 关联课程 ID
	CourseSlotID string    `json:"course_slot_id"` // 关联的课表槽位 ID（course_slot.id）
	ScheduleDate time.Time `json:"schedule_date"`  // 具体上课日期
	MissedHours  int       `json:"missed_hours"`   // 缺席课时数

	// 1 事假 / 2 公假 / 3 旷课
	AbsenceType uint8 `json:"absence_type"`

	Reason string `json:"reason"` // 缺勤事由

	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// AbsenceToRedis 写路径。这个聚合没有状态，缺勤就是一条事实。
func AbsenceToRedis(do *absence.AbsenceRecord) (*AbsenceRecord, error) {
	if do == nil {
		return nil, ErrAbsenceToRedis
	}
	return &AbsenceRecord{
		ID:           do.ID(),
		StudentID:    do.StudentID(),
		CourseID:     do.CourseID(),
		CourseSlotID: do.CourseSlotID(),
		ScheduleDate: do.ScheduleDate(),
		MissedHours:  do.MissedHours(),
		AbsenceType:  uint8(do.AbsenceType()),
		Reason:       do.Reason(),
		CreatedAt:    do.CreatedAt(),
		UpdatedAt:    do.UpdatedAt(),
	}, nil
}

// AbsenceFromRedis 读路径。
func AbsenceFromRedis(po *AbsenceRecord) (*absence.AbsenceRecord, error) {
	if po == nil {
		return nil, ErrAbsenceFromRedis
	}
	return absence.Reconstitute(
		po.ID, po.StudentID, po.CourseID, po.CourseSlotID, po.ScheduleDate,
		po.MissedHours, absence.AbsenceType(po.AbsenceType), po.Reason,
		po.CreatedAt, po.UpdatedAt,
	), nil
}
