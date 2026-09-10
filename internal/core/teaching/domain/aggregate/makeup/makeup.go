package makeup

import (
	"errors"
	"time"
)

// --- 聚合根 (Aggregate Root) ---

type StudentMakeup struct {
	id              int64
	studentID       int64     // 学员 ID
	courseID        string    // 课程 ID (原课与目标课必须同一门)
	absenceRecordID int64     // 关联需要冲销的原缺勤记录 ID

	// 原定缺席信息快照
	originalSlotID int64     // 原排课模板 ID
	originalDate   time.Time // 原缺课日期

	// 目标补课信息
	targetSlotID int64     // 目标去蹭课/补课的 CourseSlot ID
	targetDate   time.Time // 目标补课的具体日期
	makeupHours  int       // 补课冲销课时数

	status       Status    // 补课履约状态
	reviewerID   int64     // 审批人/准入确认人 (教务或目标班讲师)
	reviewRemark string    // 审批或驳回说明
	completedAt  time.Time // 实际完成核销时间
	createdAt    time.Time
	updatedAt    time.Time
}

// NewStudentMakeup 学员申请去其他课位自主补课
func NewStudentMakeup(
	id int64,
	studentID int64,
	courseID string,
	absenceRecordID int64,
	originalSlotID int64,
	originalDate time.Time,
	targetSlotID int64,
	targetDate time.Time,
	makeupHours int,
	now time.Time,
) (*StudentMakeup, error) {
	if id <= 0 || studentID <= 0 || absenceRecordID <= 0 {
		return nil, errors.New("invalid identifier")
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
		id:              id,
		studentID:       studentID,
		courseID:        courseID,
		absenceRecordID: absenceRecordID,
		originalSlotID:  originalSlotID,
		originalDate:    originalDate,
		targetSlotID:    targetSlotID,
		targetDate:      targetDate,
		makeupHours:     makeupHours,
		status:          StatusPending,
		createdAt:       now,
		updatedAt:       now,
	}, nil
}

// Reconstitute 仓储恢复聚合根
func Reconstitute(
	id int64,
	studentID int64,
	courseID string,
	absenceRecordID int64,
	originalSlotID int64,
	originalDate time.Time,
	targetSlotID int64,
	targetDate time.Time,
	makeupHours int,
	status Status,
	reviewerID int64,
	reviewRemark string,
	completedAt time.Time,
	createdAt, updatedAt time.Time,
) *StudentMakeup {
	return &StudentMakeup{
		id:              id,
		studentID:       studentID,
		courseID:        courseID,
		absenceRecordID: absenceRecordID,
		originalSlotID:  originalSlotID,
		originalDate:    originalDate,
		targetSlotID:    targetSlotID,
		targetDate:      targetDate,
		makeupHours:     makeupHours,
		status:          status,
		reviewerID:      reviewerID,
		reviewRemark:    reviewRemark,
		completedAt:     completedAt,
		createdAt:       createdAt,
		updatedAt:       updatedAt,
	}
}

// --- 核心领域行为 (Domain Behaviors) ---

// Approve 准许插班补课（目标班级有空位，教务或老师同意）
func (m *StudentMakeup) Approve(reviewerID int64, remark string, now time.Time) error {
	if m.status != StatusPending {
		return ErrAlreadyReviewed
	}
	if reviewerID <= 0 {
		return errors.New("invalid reviewer ID")
	}

	m.status = StatusApproved
	m.reviewerID = reviewerID
	m.reviewRemark = remark
	m.updatedAt = now
	return nil
}

// Reject 驳回插班补课申请（例如目标班级人数已超载、进度不匹配）
func (m *StudentMakeup) Reject(reviewerID int64, remark string, now time.Time) error {
	if m.status != StatusPending {
		return ErrAlreadyReviewed
	}
	if reviewerID <= 0 {
		return errors.New("invalid reviewer ID")
	}

	m.status = StatusRejected
	m.reviewerID = reviewerID
	m.reviewRemark = remark
	m.updatedAt = now
	return nil
}

// CompleteAttendance 补课当天现场签到核销成功
func (m *StudentMakeup) CompleteAttendance(now time.Time) error {
	if m.status != StatusApproved {
		return ErrCannotCompleteNotApp
	}

	m.status = StatusCompleted
	m.completedAt = now
	m.updatedAt = now
	return nil
}

// Cancel 学员主动取消补课预约
func (m *StudentMakeup) Cancel(operatorID int64, now time.Time) error {
	if m.status == StatusCompleted || m.status == StatusCancelled {
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

func (m *StudentMakeup) ID() int64              { return m.id }
func (m *StudentMakeup) StudentID() int64       { return m.studentID }
func (m *StudentMakeup) CourseID() string       { return m.courseID }
func (m *StudentMakeup) AbsenceRecordID() int64 { return m.absenceRecordID }
func (m *StudentMakeup) OriginalSlotID() int64  { return m.originalSlotID }
func (m *StudentMakeup) OriginalDate() time.Time { return m.originalDate }
func (m *StudentMakeup) TargetSlotID() int64    { return m.targetSlotID }
func (m *StudentMakeup) TargetDate() time.Time  { return m.targetDate }
func (m *StudentMakeup) MakeupHours() int       { return m.makeupHours }
func (m *StudentMakeup) Status() Status         { return m.status }
func (m *StudentMakeup) CompletedAt() time.Time { return m.completedAt }
