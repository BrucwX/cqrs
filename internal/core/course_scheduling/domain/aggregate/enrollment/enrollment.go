package enrollment

import (
	"errors"
	"time"
)

// --- 聚合根 (Aggregate Root) ---

type CourseEnrollment struct {
	id          int64
	studentID   int64     // 学员/员工 ID
	courseID    string    // 课程 ID
	status      Status    // 注册状态
	enrolledAt  time.Time // 报名时间
	completedAt time.Time // 结业时间
	droppedAt   time.Time // 退课时间
	updatedAt   time.Time
}

// generateID 生成唯一的 int64 ID
func generateID() int64 {
	return time.Now().UnixNano()
}

// NewCourseEnrollment 报班注册构造（报名登记：已报班 / 已缴费，尚未选课）
//
// 落在 StatusNotSelected —— 这条记录代表「选课资格」，正式选课由 Enroll 推到在读。
func NewCourseEnrollment(
	studentID int64,
	courseID string,
) (*CourseEnrollment, error) {
	if studentID <= 0 {
		return nil, errors.New("invalid student ID")
	}
	if courseID == "" {
		return nil, errors.New("course ID is required")
	}

	now := time.Now()
	return &CourseEnrollment{
		id:         generateID(),
		studentID:  studentID,
		courseID:   courseID,
		status:     StatusNotSelected,
		enrolledAt: now, // 报名登记时间
		updatedAt:  now,
	}, nil
}

// Reconstitute 仓储恢复聚合根
func Reconstitute(
	id int64,
	studentID int64,
	courseID string,
	status Status,
	enrolledAt time.Time,
	completedAt time.Time,
	droppedAt time.Time,
	updatedAt time.Time,
) *CourseEnrollment {
	return &CourseEnrollment{
		id:          id,
		studentID:   studentID,
		courseID:    courseID,
		status:      status,
		enrolledAt:  enrolledAt,
		completedAt: completedAt,
		droppedAt:   droppedAt,
		updatedAt:   updatedAt,
	}
}

// --- 核心领域行为 (Domain Behaviors) ---

// Enroll 正式选课：把「未选课的注册记录」推到在读。
//
// 不新建记录 —— 记录本身是报名登记（缴费）时落下的选课资格，
// 这里只做状态流转（准入判定由调用方在调过来之前做完）。
func (e *CourseEnrollment) Enroll(now time.Time) error {
	if e.status != StatusNotSelected {
		return ErrEnrollmentNotSelectable
	}
	e.status = StatusEnrolled
	e.updatedAt = now
	return nil
}

// Complete 结业（可由教务人员录入结业成绩、或课程整体周期结束后统一触发）
func (e *CourseEnrollment) Complete(now time.Time) error {
	if e.status != StatusEnrolled {
		return ErrEnrollmentNotActive
	}
	e.status = StatusCompleted
	e.completedAt = now
	e.updatedAt = now
	return nil
}

// Drop 申请退课
func (e *CourseEnrollment) Drop(now time.Time) error {
	if e.status == StatusCompleted {
		return ErrCannotDropCompleted
	}
	if e.status == StatusDropped {
		return ErrAlreadyDropped
	}

	e.status = StatusDropped
	e.droppedAt = now
	e.updatedAt = now
	return nil
}

// --- 只读属性访问器 (Getters) ---

func (e *CourseEnrollment) ID() int64              { return e.id }
func (e *CourseEnrollment) StudentID() int64       { return e.studentID }
func (e *CourseEnrollment) CourseID() string       { return e.courseID }
func (e *CourseEnrollment) Status() Status         { return e.status }
func (e *CourseEnrollment) EnrolledAt() time.Time  { return e.enrolledAt }
func (e *CourseEnrollment) CompletedAt() time.Time { return e.completedAt }
func (e *CourseEnrollment) DroppedAt() time.Time   { return e.droppedAt }
func (e *CourseEnrollment) UpdatedAt() time.Time   { return e.updatedAt }
func (e *CourseEnrollment) IsActive() bool         { return e.status == StatusEnrolled }
func (e *CourseEnrollment) IsCompleted() bool      { return e.status == StatusCompleted }
func (e *CourseEnrollment) IsNotSelected() bool    { return e.status == StatusNotSelected }
