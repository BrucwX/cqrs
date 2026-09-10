package courseSlotChange

import (
	"errors"
	"time"
)

// --- 聚合根 (Aggregate Root) ---

type CourseSlotChange struct {
	id           int64
	courseID     string       // 关联课程 ID
	applicantID  int64        // 发起申请人 ID（通常为原老师或教务）
	changeType   ChangeType   // 变更性质
	originalPlan OriginalPlan // 原始计划值对象
	targetPlan   TargetPlan   // 目标计划值对象
	reason       string       // 调课/代课事由
	reviewStatus ReviewStatus // 审批状态
	reviewerID   int64        // 审批人 ID
	reviewRemark string       // 审批意见
	createdAt    time.Time
	updatedAt    time.Time
}

// NewCourseSlotChange 发起临时换课申请
func NewCourseSlotChange(
	id int64,
	courseID string,
	applicantID int64,
	changeType ChangeType,
	original OriginalPlan,
	target TargetPlan,
	reason string,
	now time.Time,
) (*CourseSlotChange, error) {
	if id <= 0 {
		return nil, errors.New("invalid change ID")
	}
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
		id:           id,
		courseID:     courseID,
		applicantID:  applicantID,
		changeType:   changeType,
		originalPlan: original,
		targetPlan:   target,
		reason:       reason,
		reviewStatus: StatusPending,
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
	reviewStatus ReviewStatus,
	reviewerID int64,
	reviewRemark string,
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
		reviewStatus: reviewStatus,
		reviewerID:   reviewerID,
		reviewRemark: reviewRemark,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// --- 核心领域行为 (Domain Behaviors) ---

// Approve 审批通过：换课正式生效
func (c *CourseSlotChange) Approve(reviewerID int64, remark string, now time.Time) error {
	if c.reviewStatus != StatusPending {
		return ErrChangeAlreadyReviewed
	}
	if reviewerID <= 0 {
		return errors.New("invalid reviewer ID")
	}

	c.reviewStatus = StatusApproved
	c.reviewerID = reviewerID
	c.reviewRemark = remark
	c.updatedAt = now
	return nil
}

// Reject 审批驳回
func (c *CourseSlotChange) Reject(reviewerID int64, remark string, now time.Time) error {
	if c.reviewStatus != StatusPending {
		return ErrChangeAlreadyReviewed
	}
	if reviewerID <= 0 {
		return errors.New("invalid reviewer ID")
	}

	c.reviewStatus = StatusRejected
	c.reviewerID = reviewerID
	c.reviewRemark = remark
	c.updatedAt = now
	return nil
}

// Withdraw 申请人撤回
func (c *CourseSlotChange) Withdraw(operatorID int64, now time.Time) error {
	if c.reviewStatus != StatusPending {
		return ErrCannotCancelReviewed
	}
	if c.applicantID != operatorID {
		return errors.New("only applicant can withdraw the change request")
	}

	c.reviewStatus = StatusWithdrawn
	c.updatedAt = now
	return nil
}

// IsEffective 校验当前变更是否属于已审核通过的有效变更
func (c *CourseSlotChange) IsEffective() bool {
	return c.reviewStatus == StatusApproved
}

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
func (c *CourseSlotChange) ReviewStatus() ReviewStatus { return c.reviewStatus }
func (c *CourseSlotChange) ReviewerID() int64          { return c.reviewerID }
func (c *CourseSlotChange) ReviewRemark() string       { return c.reviewRemark }
