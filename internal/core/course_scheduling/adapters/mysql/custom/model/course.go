package model

import "time"

// Course 对应表 course，领域模型是 course.Course。
//
// 三个值对象就地打平：
//
//	Capacity         -> CapacityMax / CapacityEnrolled
//	EnrollmentWindow -> EnrollStartAt / EnrollEndAt / DropDeadline
//	CoursePeriod     -> PeriodStartAt / PeriodEndAt / TotalHours / CompletedHours
//
// CourseTypeID 指向 course_type，是本上下文里所有「按类型」规则的落点
// （讲师资质、按类型取课程）。表之间没有外键，引用只在应用里维护。
//
// 与 classroom 一样，这张表没有 created_at / updated_at。
type Course struct {
	ID           string // varchar(36) 课程 ID（uuid）
	CourseTypeID string // varchar(36) 所属课程类型（course_type.id）

	CapacityMax      int // int 总容量          <- Capacity.Max
	CapacityEnrolled int // int 已报名人数      <- Capacity.Enrolled

	EnrollStartAt time.Time // datetime 选课开始时间   <- EnrollmentWindow.StartAt
	EnrollEndAt   time.Time // datetime 选课结束时间   <- EnrollmentWindow.EndAt
	DropDeadline  time.Time // datetime 退课截止时间   <- EnrollmentWindow.DropDeadline

	PeriodStartAt  time.Time // datetime 教学周期开始    <- CoursePeriod.StartAt
	PeriodEndAt    time.Time // datetime 教学周期结束    <- CoursePeriod.EndAt
	TotalHours     int       // int      总课时          <- CoursePeriod.TotalHours
	CompletedHours int       // int      已完成课时      <- CoursePeriod.CompletedHours

	LockVersion uint64 // bigint unsigned 乐观锁版本号
}
