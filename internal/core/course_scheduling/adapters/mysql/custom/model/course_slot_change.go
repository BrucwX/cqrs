package model

import "time"

// CourseSlotChange 对应表 course_slot_change，领域模型是 courseSlotChange.CourseSlotChange。
//
// 临时换课：登记即生效，没有审批；原计划是一份快照，所以两个值对象整份打平 ——
//
//	OriginalPlan -> OriginalSlotID / OriginalDate / OriginalTeacherID /
//	                OriginalClassroomID / OriginalStartTime / OriginalEndTime
//	TargetPlan   -> TargetStartAt / TargetEndAt / TargetTeacherID / TargetClassroomID
//
// 两处「模型里就对不上」的地方（见 doc.go），这里如实建模：
//
//	a. OriginalSlotID 是 int64，但 course_slot.id 是 varchar(36) 的 uuid，
//	   这份快照指不到具体的课表模板上。
//	b. 原计划的时间是 varchar(5) 的 "09:00"，目标时间是 datetime，
//	   同一个概念在两份快照里的表示不一致。
type CourseSlotChange struct {
	ID          int64  // bigint      变更单 ID
	CourseID    string // varchar(36) 关联课程（course.id）
	ApplicantID int64  // bigint      发起申请人 ID

	// uint8 tinyint unsigned 1 改期 / 2 代课 / 3 换教室 / 4 复合变动
	ChangeType uint8

	// --- 原计划快照（OriginalPlan）---
	OriginalSlotID      int64     // bigint      原课表模板 ID（见上文 a）
	OriginalDate        time.Time // date        原定上课日期
	OriginalTeacherID   int64     // bigint      原讲师 ID
	OriginalClassroomID string    // varchar(36) 原教室 ID
	OriginalStartTime   string    // varchar(5)  原开始时间（"09:00"）
	OriginalEndTime     string    // varchar(5)  原结束时间（"09:00"）

	// --- 目标计划（TargetPlan）---
	TargetStartAt     time.Time // datetime    调整后开始时间
	TargetEndAt       time.Time // datetime    调整后结束时间（与开始同日）
	TargetTeacherID   int64     // bigint      实际授课讲师 ID
	TargetClassroomID string    // varchar(36) 实际使用教室 ID

	Reason string // varchar(255) 调课/代课事由

	CreatedAt   time.Time // datetime 创建时间
	UpdatedAt   time.Time // datetime 更新时间
	LockVersion uint64    // bigint unsigned 乐观锁版本号
}
