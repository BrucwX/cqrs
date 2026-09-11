package makeup

import (
	"errors"
	"time"
)

// --- 聚合根 (Aggregate Root) ---

// StudentMakeup 补课预约：约好了到时候去上，不需要审批。
type StudentMakeup struct {
	id        int64
	studentID int64  // 学员 ID
	courseID  string // 课程 ID（补课不跨课程）

	// 原定缺席信息快照
	originalSlotID int64     // 原排课模板 ID
	originalDate   time.Time // 原缺课日期

	// 目标补课信息
	targetSlotID int64     // 目标去蹭课/补课的 CourseSlot ID
	targetDate   time.Time // 目标补课的具体日期
	makeupHours  int       // 补课课时数

	status      Status    // 补课状态
	completedAt time.Time // 实际完成核销时间
	createdAt   time.Time
	updatedAt   time.Time
}

// generateID 生成唯一的 int64 ID
func generateID() int64 {
	return time.Now().UnixNano()
}

// NewStudentMakeup 学员预约补课
func NewStudentMakeup(
	studentID int64,
	courseID string,
	originalSlotID int64,
	originalDate time.Time,
	targetSlotID int64,
	targetDate time.Time,
	makeupHours int,
	now time.Time,
) (*StudentMakeup, error) {
	if studentID <= 0 {
		return nil, errors.New("invalid student ID")
	}
	if courseID == "" {
		return nil, errors.New("course ID is required")
	}
	if targetSlotID <= 0 {
		return nil, errors.New("target course slot ID is required")
	}
	if makeupHours <= 0 {
		return nil, errors.New("makeup hours must be greater than zero")
	}

	return &StudentMakeup{
		id:             generateID(),
		studentID:      studentID,
		courseID:       courseID,
		originalSlotID: originalSlotID,
		originalDate:   originalDate,
		targetSlotID:   targetSlotID,
		targetDate:     targetDate,
		makeupHours:    makeupHours,
		status:         StatusBooked,
		createdAt:      now,
		updatedAt:      now,
	}, nil
}

// Reconstitute 仓储恢复聚合根
func Reconstitute(
	id int64,
	studentID int64,
	courseID string,
	originalSlotID int64,
	originalDate time.Time,
	targetSlotID int64,
	targetDate time.Time,
	makeupHours int,
	status Status,
	completedAt time.Time,
	createdAt, updatedAt time.Time,
) *StudentMakeup {
	return &StudentMakeup{
		id:             id,
		studentID:      studentID,
		courseID:       courseID,
		originalSlotID: originalSlotID,
		originalDate:   originalDate,
		targetSlotID:   targetSlotID,
		targetDate:     targetDate,
		makeupHours:    makeupHours,
		status:         status,
		completedAt:    completedAt,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

// --- 核心领域行为 (Domain Behaviors) ---

// CompleteAttendance 补课当天现场签到核销成功
func (m *StudentMakeup) CompleteAttendance(now time.Time) error {
	if m.status != StatusBooked {
		return ErrAlreadyFinalized
	}

	m.status = StatusCompleted
	m.completedAt = now
	m.updatedAt = now
	return nil
}

// Cancel 学员主动取消补课预约
func (m *StudentMakeup) Cancel(operatorID int64, now time.Time) error {
	if m.status != StatusBooked {
		return ErrAlreadyFinalized
	}
	if m.studentID != operatorID {
		return errors.New("only applicant can cancel this makeup reservation")
	}

	m.status = StatusCancelled
	m.updatedAt = now
	return nil
}

// --- 只读属性访问器 (Getters) ---

func (m *StudentMakeup) ID() int64               { return m.id }
func (m *StudentMakeup) StudentID() int64        { return m.studentID }
func (m *StudentMakeup) CourseID() string        { return m.courseID }
func (m *StudentMakeup) OriginalSlotID() int64   { return m.originalSlotID }
func (m *StudentMakeup) OriginalDate() time.Time { return m.originalDate }
func (m *StudentMakeup) TargetSlotID() int64     { return m.targetSlotID }
func (m *StudentMakeup) TargetDate() time.Time   { return m.targetDate }
func (m *StudentMakeup) MakeupHours() int        { return m.makeupHours }
func (m *StudentMakeup) Status() Status          { return m.status }
func (m *StudentMakeup) CompletedAt() time.Time  { return m.completedAt }
