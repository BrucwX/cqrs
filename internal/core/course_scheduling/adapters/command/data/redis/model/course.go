package model

import (
	"fmt"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
)

// Course 是 Redis 快照（对应 course），领域模型是 course.Course。
//
// 三个值对象就地打平：
//
//	Capacity         -> CapacityMax / CapacityEnrolled
//	EnrollmentWindow -> EnrollStartAt / EnrollEndAt / DropDeadline
//	CoursePeriod     -> PeriodStartAt / PeriodEndAt / TotalHours / CompletedHours
//
// CourseTypeID 指向 course_type，是本上下文里所有「按类型」规则的落点
// （讲师资质、按类型取课程）。聚合之间不做级联，引用只在应用里维护。
//
// 与 classroom 一样，快照里没有 created_at / updated_at。
type Course struct {
	ID           string `json:"id"`             // 课程 ID（uuid）
	CourseTypeID string `json:"course_type_id"` // 所属课程类型（course_type.id）

	CapacityMax      int `json:"capacity_max"`      // 总容量          <- Capacity.Max
	CapacityEnrolled int `json:"capacity_enrolled"` // 已报名人数      <- Capacity.Enrolled

	EnrollStartAt time.Time `json:"enroll_start_at"` // 选课开始时间   <- EnrollmentWindow.StartAt
	EnrollEndAt   time.Time `json:"enroll_end_at"`   // 选课结束时间   <- EnrollmentWindow.EndAt
	DropDeadline  time.Time `json:"drop_deadline"`   // 退课截止时间   <- EnrollmentWindow.DropDeadline

	PeriodStartAt  time.Time `json:"period_start_at"` // 教学周期开始    <- CoursePeriod.StartAt
	PeriodEndAt    time.Time `json:"period_end_at"`   // 教学周期结束    <- CoursePeriod.EndAt
	TotalHours     int       `json:"total_hours"`     // 总课时          <- CoursePeriod.TotalHours
	CompletedHours int       `json:"completed_hours"` // 已完成课时      <- CoursePeriod.CompletedHours

}

// CourseToRedis 写路径：Capacity / EnrollmentWindow / CoursePeriod 三个值对象打平。
func CourseToRedis(do *course.Course) (*Course, error) {
	if do == nil {
		return nil, ErrCourseToRedis
	}
	capacity := do.Capacity()
	window := do.Enrollment()
	period := do.Period()
	return &Course{
		ID:               do.ID(),
		CourseTypeID:     do.CourseTypeID(),
		CapacityMax:      capacity.Max(),
		CapacityEnrolled: capacity.Enrolled(),
		EnrollStartAt:    window.StartAt(),
		EnrollEndAt:      window.EndAt(),
		DropDeadline:     window.DropDeadline(),
		PeriodStartAt:    period.StartAt(),
		PeriodEndAt:      period.EndAt(),
		TotalHours:       period.TotalHours(),
		CompletedHours:   period.CompletedHours(),
	}, nil
}

// CourseFromRedis 读路径：容量与课时都走值对象构造函数，越界数据在这里报错。
func CourseFromRedis(po *Course) (*course.Course, error) {
	if po == nil {
		return nil, ErrCourseFromRedis
	}
	capacity, err := course.NewCapacity(po.CapacityMax, po.CapacityEnrolled)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCourseFromRedis, err)
	}
	period, err := course.NewCoursePeriod(
		po.PeriodStartAt, po.PeriodEndAt, po.TotalHours, po.CompletedHours,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCourseFromRedis, err)
	}
	window := course.NewEnrollmentWindow(po.EnrollStartAt, po.EnrollEndAt, po.DropDeadline)

	return course.Reconstitute(po.ID, po.CourseTypeID, capacity, window, period), nil
}
