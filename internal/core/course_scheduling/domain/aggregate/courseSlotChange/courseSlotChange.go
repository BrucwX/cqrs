package courseSlotChange

import (
	"errors"
	"time"
)

// --- 聚合根 (Aggregate Root) ---

// CourseSlotChange 课表变更记录：把某节课换到别的时间/别的人。
//
// 登记即生效，没有审批环节。
type CourseSlotChange struct {
	id           int64
	courseID     string       // 关联课程 ID
	applicantID  int64        // 发起申请人 ID（通常为原老师或教务）
	changeType   ChangeType   // 变更性质
	originalPlan OriginalPlan // 原始计划值对象
	targetPlan   TargetPlan   // 目标计划值对象
	reason       string       // 调课/代课事由
	createdAt    time.Time
	updatedAt    time.Time
}

// generateID 生成唯一的 int64 ID
func generateID() int64 {
	return time.Now().UnixNano()
}

// NewCourseSlotChange 登记一次临时换课
func NewCourseSlotChange(
	courseID string,
	applicantID int64,
	changeType ChangeType,
	original OriginalPlan,
	target TargetPlan,
	reason string,
	now time.Time,
) (*CourseSlotChange, error) {
	if courseID == "" {
		return nil, errors.New("course ID is required")
	}
	if applicantID <= 0 {
		return nil, errors.New("applicant ID is required")
	}
	if reason == "" {
		return nil, errors.New("change reason is required")
	}
	if target.targetStartAt.Before(now) {
		return nil, ErrTargetDateInPast
	}

	return &CourseSlotChange{
		id:           generateID(),
		courseID:     courseID,
		applicantID:  applicantID,
		changeType:   changeType,
		originalPlan: original,
		targetPlan:   target,
		reason:       reason,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// Reconstitute 仓储恢复聚合根
func Reconstitute(
	id int64,
	courseID string,
	applicantID int64,
	changeType ChangeType,
	original OriginalPlan,
	target TargetPlan,
	reason string,
	createdAt, updatedAt time.Time,
) *CourseSlotChange {
	return &CourseSlotChange{
		id:           id,
		courseID:     courseID,
		applicantID:  applicantID,
		changeType:   changeType,
		originalPlan: original,
		targetPlan:   target,
		reason:       reason,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// --- 核心领域行为 (Domain Behaviors) ---

// IsSubstituteTeacher 是否涉及代课老师变动
func (c *CourseSlotChange) IsSubstituteTeacher() bool {
	return c.originalPlan.teacherID != c.targetPlan.teacherID
}

// --- 只读属性访问器 (Getters) ---

func (c *CourseSlotChange) ID() int64                  { return c.id }
func (c *CourseSlotChange) CourseID() string           { return c.courseID }
func (c *CourseSlotChange) ApplicantID() int64         { return c.applicantID }
func (c *CourseSlotChange) ChangeType() ChangeType     { return c.changeType }
func (c *CourseSlotChange) OriginalPlan() OriginalPlan { return c.originalPlan }
func (c *CourseSlotChange) TargetPlan() TargetPlan     { return c.targetPlan }
func (c *CourseSlotChange) Reason() string             { return c.reason }
