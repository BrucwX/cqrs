package model

import (
	"fmt"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
)

// CourseSlotChange 是 Redis 快照（对应 course_slot_change），领域模型是 courseSlotChange.CourseSlotChange。
//
// 临时换课：登记即生效，没有审批；原计划是一份快照，所以两个值对象整份打平 ——
//
//	OriginalPlan -> OriginalSlotID / OriginalDate / OriginalTeacherID /
//	                OriginalClassroomID / OriginalStartTime / OriginalEndTime
//	TargetPlan   -> TargetStartAt / TargetEndAt / TargetTeacherID / TargetClassroomID
//
// 一处「模型里就对不上」的地方（见 doc.go），这里如实建模：
//
//	a. 原计划的时间是 varchar(5) 的 "09:00"，目标时间是 datetime，
//	   同一个概念在两份快照里的表示不一致。
type CourseSlotChange struct {
	ID          int64  `json:"id"`           // 变更单 ID
	CourseID    string `json:"course_id"`    // 关联课程（course.id）
	ApplicantID int64  `json:"applicant_id"` // 发起申请人 ID

	// 1 改期 / 2 代课 / 3 换教室 / 4 复合变动
	ChangeType uint8 `json:"change_type"`

	// --- 原计划快照（OriginalPlan）---
	OriginalSlotID      string    `json:"original_slot_id"`      // 原课表槽位 ID（course_slot.id）
	OriginalDate        time.Time `json:"original_date"`         // 原定上课日期
	OriginalTeacherID   int64     `json:"original_teacher_id"`   // 原讲师 ID
	OriginalClassroomID string    `json:"original_classroom_id"` // 原教室 ID
	OriginalStartTime   string    `json:"original_start_time"`   // 原开始时间（"09:00"）
	OriginalEndTime     string    `json:"original_end_time"`     // 原结束时间（"09:00"）

	// --- 目标计划（TargetPlan）---
	TargetStartAt     time.Time `json:"target_start_at"`     // 调整后开始时间
	TargetEndAt       time.Time `json:"target_end_at"`       // 调整后结束时间（与开始同日）
	TargetTeacherID   int64     `json:"target_teacher_id"`   // 实际授课讲师 ID
	TargetClassroomID string    `json:"target_classroom_id"` // 实际使用教室 ID

	Reason string `json:"reason"` // 调课/代课事由

	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// CourseSlotChangeToRedis 写路径：OriginalPlan 打平成 original_* 六个字段，
// TargetPlan 打平成 target_* 四个字段。
func CourseSlotChangeToRedis(do *courseSlotChange.CourseSlotChange) (*CourseSlotChange, error) {
	if do == nil {
		return nil, ErrCourseSlotChangeToRedis
	}
	original := do.OriginalPlan()
	target := do.TargetPlan()
	return &CourseSlotChange{
		ID:          do.ID(),
		CourseID:    do.CourseID(),
		ApplicantID: do.ApplicantID(),
		ChangeType:  uint8(do.ChangeType()),

		OriginalSlotID:      original.SlotID(),
		OriginalDate:        original.Date(),
		OriginalTeacherID:   original.TeacherID(),
		OriginalClassroomID: original.ClassroomID(),
		OriginalStartTime:   original.StartTimeStr(),
		OriginalEndTime:     original.EndTimeStr(),

		TargetStartAt:     target.TargetStartAt(),
		TargetEndAt:       target.TargetEndAt(),
		TargetTeacherID:   target.TeacherID(),
		TargetClassroomID: target.ClassroomID(),

		Reason:    do.Reason(),
		CreatedAt: do.CreatedAt(),
		UpdatedAt: do.UpdatedAt(),
	}, nil
}

// CourseSlotChangeFromRedis 读路径：目标计划要走 NewTargetPlan，跨天或时间倒挂会报错。
//
// 原计划的 original_start_time / original_end_time 就是字符串，
// 和值对象一一对应，直接带过去。
func CourseSlotChangeFromRedis(po *CourseSlotChange) (*courseSlotChange.CourseSlotChange, error) {
	if po == nil {
		return nil, ErrCourseSlotChangeFromRedis
	}
	original := courseSlotChange.NewOriginalPlan(
		po.OriginalSlotID, po.OriginalDate, po.OriginalTeacherID,
		po.OriginalClassroomID, po.OriginalStartTime, po.OriginalEndTime,
	)
	target, err := courseSlotChange.NewTargetPlan(
		po.TargetStartAt, po.TargetEndAt, po.TargetTeacherID, po.TargetClassroomID,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCourseSlotChangeFromRedis, err)
	}
	return courseSlotChange.Reconstitute(
		po.ID, po.CourseID, po.ApplicantID, courseSlotChange.ChangeType(po.ChangeType),
		original, target, po.Reason, po.CreatedAt, po.UpdatedAt,
	), nil
}
