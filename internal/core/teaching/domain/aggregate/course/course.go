package course

import "time"

// --- 聚合根 (Aggregate Root) ---

type Course struct {
	id         string
	isRetake   bool
	capacity   Capacity
	enrollment EnrollmentWindow
	period     CoursePeriod
}

// 工厂方法创建 Course 聚合根
func NewCourse(
	id string,
	isRetake bool,
	capacity Capacity,
	enrollment EnrollmentWindow,
	period CoursePeriod,
) *Course {
	return &Course{
		id:         id,
		isRetake:   isRetake,
		capacity:   capacity,
		enrollment: enrollment,
		period:     period,
	}
}

// --- 核心业务行为 (Domain Behaviors) ---

// Enroll 选课行为：校验时间窗口与库存容量
func (c *Course) Enroll(now time.Time) error {
	if !c.enrollment.CanEnroll(now) {
		return ErrNotInEnrollmentStage
	}
	if c.capacity.IsFull() {
		return ErrCourseFull
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

// --- Getter（仅读，无 Setter）---

func (c *Course) ID() string                   { return c.id }
func (c *Course) IsRetake() bool               { return c.isRetake }
func (c *Course) Capacity() Capacity           { return c.capacity }
func (c *Course) Enrollment() EnrollmentWindow { return c.enrollment }
func (c *Course) Period() CoursePeriod         { return c.period }
