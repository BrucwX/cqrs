package absence

import (
	"errors"
	"time"
)

// --- 聚合根 (Aggregate Root) ---

type AbsenceRecord struct {
	id           int64
	studentID    int64        // 学员/员工 ID
	courseID     string       // 关联课程 ID
	courseSlotID int64        // 关联的具体排课模板槽位 ID
	scheduleDate time.Time    // 具体上课日期（年月日）
	missedHours  int          // 缺席课时数
	absenceType  AbsenceType  // 缺勤类型
	reason       string       // 请假/缺席事由
	reviewStatus ReviewStatus // 审核状态
	reviewerID   int64        // 审核人 ID (讲师/教务)
	reviewRemark string       // 审批意见
	isRectified  bool         // 是否已通过"补卡"或"补训"冲销
	createdAt    time.Time
	updatedAt    time.Time
}

// NewLeaveRequest 学员发起请假申请（待审批）
func NewLeaveRequest(
	id int64,
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
		id:           id,
		studentID:    studentID,
		courseID:     courseID,
		courseSlotID: courseSlotID,
		scheduleDate: scheduleDate,
		missedHours:  missedHours,
		absenceType:  absenceType,
		reason:       reason,
		reviewStatus: StatusPending,
		isRectified:  false,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// NewUnexcusedRecord 考勤点名录入旷课（无需审核，直接生效）
func NewUnexcusedRecord(
	id int64,
	studentID int64,
	courseID string,
	courseSlotID int64,
	scheduleDate time.Time,
	missedHours int,
	reason string,
) (*AbsenceRecord, error) {
	if studentID <= 0 || courseSlotID <= 0 {
		return nil, errors.New("invalid student ID or courseSlot ID")
	}
	if missedHours <= 0 {
		return nil, ErrInvalidHours
	}

	now := time.Now()
	return &AbsenceRecord{
		id:           id,
		studentID:    studentID,
		courseID:     courseID,
		courseSlotID: courseSlotID,
		scheduleDate: scheduleDate,
		missedHours:  missedHours,
		absenceType:  TypeUnexcused,
		reason:       reason,
		reviewStatus: StatusApproved,
		isRectified:  false,
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
	reviewStatus ReviewStatus,
	reviewerID int64,
	reviewRemark string,
	isRectified bool,
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
		reviewStatus: reviewStatus,
		reviewerID:   reviewerID,
		reviewRemark: reviewRemark,
		isRectified:  isRectified,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// Approve 批准请假
func (a *AbsenceRecord) Approve(reviewerID int64, remark string) error {
	if a.reviewStatus != StatusPending {
		return ErrAlreadyReviewed
	}
	a.reviewStatus = StatusApproved
	a.reviewerID = reviewerID
	a.reviewRemark = remark
	a.updatedAt = time.Now()
	return nil
}

// Reject 驳回请假，自动定性为旷课
func (a *AbsenceRecord) Reject(reviewerID int64, remark string) error {
	if a.reviewStatus != StatusPending {
		return ErrAlreadyReviewed
	}
	a.reviewStatus = StatusRejected
	a.absenceType = TypeUnexcused
	a.reviewerID = reviewerID
	a.reviewRemark = remark
	a.updatedAt = time.Now()
	return nil
}

// MarkAsRectified 被补卡或补课成功冲销
func (a *AbsenceRecord) MarkAsRectified() error {
	if a.isRectified {
		return ErrAlreadyRectified
	}
	a.isRectified = true
	a.updatedAt = time.Now()
	return nil
}

// Getters
func (a *AbsenceRecord) ID() int64                  { return a.id }
func (a *AbsenceRecord) StudentID() int64           { return a.studentID }
func (a *AbsenceRecord) CourseID() string           { return a.courseID }
func (a *AbsenceRecord) CourseSlotID() int64        { return a.courseSlotID }
func (a *AbsenceRecord) ScheduleDate() time.Time    { return a.scheduleDate }
func (a *AbsenceRecord) MissedHours() int           { return a.missedHours }
func (a *AbsenceRecord) AbsenceType() AbsenceType   { return a.absenceType }
func (a *AbsenceRecord) ReviewStatus() ReviewStatus { return a.reviewStatus }
func (a *AbsenceRecord) IsRectified() bool          { return a.isRectified }
