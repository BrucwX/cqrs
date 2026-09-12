package absence

import (
	"errors"
	"time"
)

// --- 聚合根 (Aggregate Root) ---

// AbsenceRecord 缺勤记录：学员某节课没来。
//
// 就是一条事实，没有审批、也没有「是否已补卡/已冲销」这类状态。
type AbsenceRecord struct {
	id           int64
	studentID    int64       // 学员/员工 ID
	courseID     string      // 关联课程 ID
	courseSlotID int64       // 关联的具体排课模板槽位 ID
	scheduleDate time.Time   // 具体上课日期（年月日）
	missedHours  int         // 缺席课时数
	absenceType  AbsenceType // 缺勤类型（事假 / 公假 / 旷课）
	reason       string      // 缺勤事由
	createdAt    time.Time
	updatedAt    time.Time
}

// generateID 生成唯一的 int64 ID
func generateID() int64 {
	return time.Now().UnixNano()
}

// NewAbsenceRecord 登记一条缺勤记录
func NewAbsenceRecord(
	studentID int64,
	courseID string,
	courseSlotID int64,
	scheduleDate time.Time,
	missedHours int,
	absenceType AbsenceType,
	reason string,
) (*AbsenceRecord, error) {
	if studentID <= 0 || courseSlotID <= 0 {
		return nil, errors.New("invalid student ID or courseSlot ID")
	}
	if courseID == "" {
		return nil, errors.New("course ID is required")
	}
	if missedHours <= 0 {
		return nil, ErrInvalidHours
	}

	now := time.Now()
	return &AbsenceRecord{
		id:           generateID(),
		studentID:    studentID,
		courseID:     courseID,
		courseSlotID: courseSlotID,
		scheduleDate: scheduleDate,
		missedHours:  missedHours,
		absenceType:  absenceType,
		reason:       reason,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// Reconstitute 仓储层还原
func Reconstitute(
	id int64,
	studentID int64,
	courseID string,
	courseSlotID int64,
	scheduleDate time.Time,
	missedHours int,
	absenceType AbsenceType,
	reason string,
	createdAt, updatedAt time.Time,
) *AbsenceRecord {
	return &AbsenceRecord{
		id:           id,
		studentID:    studentID,
		courseID:     courseID,
		courseSlotID: courseSlotID,
		scheduleDate: scheduleDate,
		missedHours:  missedHours,
		absenceType:  absenceType,
		reason:       reason,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// Getters
func (a *AbsenceRecord) ID() int64                { return a.id }
func (a *AbsenceRecord) StudentID() int64         { return a.studentID }
func (a *AbsenceRecord) CourseID() string         { return a.courseID }
func (a *AbsenceRecord) CourseSlotID() int64      { return a.courseSlotID }
func (a *AbsenceRecord) ScheduleDate() time.Time  { return a.scheduleDate }
func (a *AbsenceRecord) MissedHours() int         { return a.missedHours }
func (a *AbsenceRecord) AbsenceType() AbsenceType { return a.absenceType }
func (a *AbsenceRecord) Reason() string           { return a.reason }
func (a *AbsenceRecord) CreatedAt() time.Time     { return a.createdAt }
func (a *AbsenceRecord) UpdatedAt() time.Time     { return a.updatedAt }
