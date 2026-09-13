package model

import (
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
)

// AbsenceRecord 对应表 absence_record，领域模型是 absence.AbsenceRecord。
//
// 一条事实：没有审批、也没有「已补卡」这类状态，所以表里没有任何状态列。
type AbsenceRecord struct {
	ID           int64     // bigint      缺勤记录 ID
	StudentID    int64     // bigint      学员/员工 ID
	CourseID     string    // varchar(36) 关联课程 ID
	CourseSlotID string    // varchar(36) 关联的课表槽位 ID（course_slot.id）
	ScheduleDate time.Time // date        具体上课日期
	MissedHours  int       // int         缺席课时数

	// uint8 tinyint unsigned 1 事假 / 2 公假 / 3 旷课
	AbsenceType uint8

	Reason string // varchar(255) 缺勤事由

	CreatedAt   time.Time // datetime 创建时间
	UpdatedAt   time.Time // datetime 更新时间
	LockVersion uint64    // bigint unsigned 乐观锁版本号
}

// AbsenceToPO 写路径。这张表没有状态列，缺勤就是一条事实。
func AbsenceToPO(do *absence.AbsenceRecord) (*AbsenceRecord, error) {
	if do == nil {
		return nil, ErrAbsenceDOToPO
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

// AbsenceToDO 读路径。
func AbsenceToDO(po *AbsenceRecord) (*absence.AbsenceRecord, error) {
	if po == nil {
		return nil, ErrAbsencePOToDO
	}
	return absence.Reconstitute(
		po.ID, po.StudentID, po.CourseID, po.CourseSlotID, po.ScheduleDate,
		po.MissedHours, absence.AbsenceType(po.AbsenceType), po.Reason,
		po.CreatedAt, po.UpdatedAt,
	), nil
}
