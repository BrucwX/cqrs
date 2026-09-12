package model

import (
	"database/sql"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
)

// StudentMakeup 对应表 student_makeup，领域模型是 makeup.StudentMakeup。
//
// 预约即生效，补课不跨课程，所以 CourseID 与两端槽位都在同一门课内。
// 原缺课信息与目标补课信息都是快照，各自成对打平：
//
//	OriginalSlotID / OriginalDate -> 原本缺席的那次课
//	TargetSlotID   / TargetDate   -> 去补的那次课
//
// CompletedAt 为 NULL 表示还没现场核销（领域模型里是零值时间）。
type StudentMakeup struct {
	ID        int64  // bigint      补课记录 ID
	StudentID int64  // bigint      学员 ID
	CourseID  string // varchar(36) 课程 ID（补课不跨课程）

	OriginalSlotID int64     // bigint 原排课槽位 ID
	OriginalDate   time.Time // date   原缺课日期

	TargetSlotID int64     // bigint 目标补课槽位 ID
	TargetDate   time.Time // date   目标补课日期

	MakeupHours int // int 补课课时数

	// uint8 tinyint unsigned 1 已预约 / 2 已补课 / 3 已取消
	Status uint8

	CompletedAt sql.NullTime // datetime 现场核销时间（未核销为 NULL）

	CreatedAt   time.Time // datetime 创建时间
	UpdatedAt   time.Time // datetime 更新时间
	LockVersion uint64    // bigint unsigned 乐观锁版本号
}

// MakeupToPO 写路径：CompletedAt 零值表示还没现场核销，写进可空列就是 NULL。
func MakeupToPO(do *makeup.StudentMakeup) *StudentMakeup {
	return &StudentMakeup{
		ID:             do.ID(),
		StudentID:      do.StudentID(),
		CourseID:       do.CourseID(),
		OriginalSlotID: do.OriginalSlotID(),
		OriginalDate:   do.OriginalDate(),
		TargetSlotID:   do.TargetSlotID(),
		TargetDate:     do.TargetDate(),
		MakeupHours:    do.MakeupHours(),
		Status:         uint8(do.Status()),
		CompletedAt:    nullTime(do.CompletedAt()),
		CreatedAt:      do.CreatedAt(),
		UpdatedAt:      do.UpdatedAt(),
	}
}

// MakeupToDO 读路径：NULL 还原成零值时间。
func MakeupToDO(po *StudentMakeup) (*makeup.StudentMakeup, error) {
	return makeup.Reconstitute(
		po.ID, po.StudentID, po.CourseID,
		po.OriginalSlotID, po.OriginalDate, po.TargetSlotID, po.TargetDate,
		po.MakeupHours, makeup.Status(po.Status), timeValue(po.CompletedAt),
		po.CreatedAt, po.UpdatedAt,
	), nil
}
