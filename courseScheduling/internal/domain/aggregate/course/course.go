package course

import (
	"time"

	"github.com/google/uuid"
)

// --- 聚合根 (Aggregate Root) ---

type Course struct {
	id           string
	courseTypeID string // 所属课程类型（courseType.CourseType）
	capacity     Capacity
	enrollment   EnrollmentWindow
	period       CoursePeriod
}

// NewCourse 创建课程（ID 由聚合自己生成）
//
// courseTypeID / capacity / enrollment / period 必填（指针为 nil 表示没给）。
func NewCourse(
	courseTypeID *string,
	capacity *Capacity,
	enrollment *EnrollmentWindow,
	period *CoursePeriod,
) (*Course, error) {
	if courseTypeID == nil || *courseTypeID == "" {
		return nil, ErrCourseTypeRequired
	}
	if capacity == nil {
		return nil, ErrCapacityRequired
	}
	if enrollment == nil {
		return nil, ErrEnrollmentRequired
	}
	if period == nil {
		return nil, ErrPeriodRequired
	}

	return &Course{
		id:           uuid.New().String(),
		courseTypeID: *courseTypeID,
		capacity:     *capacity,
		enrollment:   *enrollment,
		period:       *period,
	}, nil
}

// Reconstitute 仓储/种子数据恢复：ID 由外部给定。
func Reconstitute(
	id string,
	courseTypeID string,
	capacity Capacity,
	enrollment EnrollmentWindow,
	period CoursePeriod,
) *Course {
	return &Course{
		id:           id,
		courseTypeID: courseTypeID,
		capacity:     capacity,
		enrollment:   enrollment,
		period:       period,
	}
}

// --- 核心业务行为 (Domain Behaviors) ---

// CanEnroll 选课准入检查（只读，不改状态）：是否在选课窗口内、是否还有名额。
func (c *Course) CanEnroll(now time.Time) error {
	if !c.enrollment.CanEnroll(now) {
		return ErrNotInEnrollmentStage
	}
	if c.capacity.IsFull() {
		return ErrCourseFull
	}
	return nil
}

// Enroll 选课行为：校验时间窗口与库存容量
func (c *Course) Enroll(now time.Time) error {
	if err := c.CanEnroll(now); err != nil {
		return err
	}
	c.capacity.enrolled++
	return nil
}

// Drop 退课行为：校验退课截止时间
func (c *Course) Drop(now time.Time) error {
	if !c.enrollment.CanDrop(now) {
		return ErrDropDeadlinePassed
	}
	if c.capacity.enrolled > 0 {
		c.capacity.enrolled--
	}
	return nil
}

// RecordCompletedHours 登记课时完成进度
func (c *Course) RecordCompletedHours(hours int) error {
	newCompleted := c.period.completedHours + hours
	if newCompleted > c.period.totalHours {
		return ErrInvalidHours
	}
	c.period.completedHours = newCompleted
	return nil
}

// Update 按非 nil 的字段更新课程，nil 的字段保持原值。
func (c *Course) Update(
	courseTypeID *string,
	capacity *Capacity,
	enrollment *EnrollmentWindow,
	period *CoursePeriod,
) error {
	if courseTypeID != nil {
		if *courseTypeID == "" {
			return ErrCourseTypeRequired
		}
		c.courseTypeID = *courseTypeID
	}
	if capacity != nil {
		c.capacity = *capacity
	}
	if enrollment != nil {
		c.enrollment = *enrollment
	}
	if period != nil {
		c.period = *period
	}
	return nil
}

// --- Getter（仅读，无 Setter）---

func (c *Course) ID() string                   { return c.id }
func (c *Course) CourseTypeID() string         { return c.courseTypeID }
func (c *Course) Capacity() Capacity           { return c.capacity }
func (c *Course) Enrollment() EnrollmentWindow { return c.enrollment }
func (c *Course) Period() CoursePeriod         { return c.period }
