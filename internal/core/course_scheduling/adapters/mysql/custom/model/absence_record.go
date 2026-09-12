package model

import "time"

// AbsenceRecord 对应表 absence_record，领域模型是 absence.AbsenceRecord。
//
// 一条事实：没有审批、也没有「已补卡」这类状态，所以表里没有任何状态列。
//
// CourseSlotID 与 course_slot.id 类型对不上（int64 vs varchar(36)），
// 见 doc.go 的说明；这里按 Go 类型如实建模。
type AbsenceRecord struct {
	ID           int64     // bigint      缺勤记录 ID
	StudentID    int64     // bigint      学员/员工 ID
	CourseID     string    // varchar(36) 关联课程 ID
	CourseSlotID int64     // bigint      关联的排课槽位 ID
	ScheduleDate time.Time // date        具体上课日期
	MissedHours  int       // int         缺席课时数

	// uint8 tinyint unsigned 1 事假 / 2 公假 / 3 旷课
	AbsenceType uint8

	Reason string // varchar(255) 缺勤事由

	CreatedAt   time.Time // datetime 创建时间
	UpdatedAt   time.Time // datetime 更新时间
	LockVersion uint64    // bigint unsigned 乐观锁版本号
}
