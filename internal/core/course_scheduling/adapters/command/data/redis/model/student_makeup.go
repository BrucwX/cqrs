package model

import (
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
)

// StudentMakeup 是 Redis 快照（对应 student_makeup），领域模型是 makeup.StudentMakeup。
//
// 预约即生效，补课不跨课程，所以 CourseID 与两端槽位都在同一门课内。
// 原缺课信息与目标补课信息都是快照，各自成对打平：
//
//	OriginalSlotID / OriginalDate -> 原本缺席的那次课
//	TargetSlotID   / TargetDate   -> 去补的那次课
//
// CompletedAt 为 null 表示还没现场核销（领域模型里是零值时间）。
type StudentMakeup struct {
	ID        int64  `json:"id"`         // 补课记录 ID
	StudentID int64  `json:"student_id"` // 学员 ID
	CourseID  string `json:"course_id"`  // 课程 ID（补课不跨课程）

	OriginalSlotID string    `json:"original_slot_id"` // 原本缺席的那节 CourseSlot ID
	OriginalDate   time.Time `json:"original_date"`    // 原缺课日期

	TargetSlotID string    `json:"target_slot_id"` // 目标补课的那节 CourseSlot ID
	TargetDate   time.Time `json:"target_date"`    // 目标补课日期

	MakeupHours int `json:"makeup_hours"` // 补课课时数

	// 1 已预约 / 2 已补课 / 3 已取消
	Status uint8 `json:"status"`

	CompletedAt *time.Time `json:"completed_at"` // 现场核销时间（未核销为 null）

	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// MakeupToRedis 写路径：CompletedAt 零值表示还没现场核销，写成 JSON 就是 null。
func MakeupToRedis(do *makeup.StudentMakeup) (*StudentMakeup, error) {
	if do == nil {
		return nil, ErrMakeupToRedis
	}
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
		CompletedAt:    timePtr(do.CompletedAt()),
		CreatedAt:      do.CreatedAt(),
		UpdatedAt:      do.UpdatedAt(),
	}, nil
}

// MakeupFromRedis 读路径：null 还原成零值时间。
func MakeupFromRedis(po *StudentMakeup) (*makeup.StudentMakeup, error) {
	if po == nil {
		return nil, ErrMakeupFromRedis
	}
	return makeup.Reconstitute(
		po.ID, po.StudentID, po.CourseID,
		po.OriginalSlotID, po.OriginalDate, po.TargetSlotID, po.TargetDate,
		po.MakeupHours, makeup.Status(po.Status), timeValue(po.CompletedAt),
		po.CreatedAt, po.UpdatedAt,
	), nil
}
