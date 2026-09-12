package model

import (
	"database/sql"
	"time"
)

// CourseEnrollment 对应表 course_enrollment，领域模型是 enrollment.CourseEnrollment。
//
// 同一个学员可以重复报同一门课（退课后重选），所以**没有唯一键**，
// 履约过程由 Status 状态机表达（在读 / 已结业 / 已退课）。
//
// 两个可空时间列在领域模型里是零值时间，在这里就是 NULL：
//
//	CompletedAt NULL = 未结业
//	DroppedAt   NULL = 未退课
type CourseEnrollment struct {
	ID        int64  // bigint      报名记录 ID
	StudentID int64  // bigint      学员 ID
	CourseID  string // varchar(36) 课程 ID（course.id）

	// uint8 tinyint unsigned 1 在读 / 2 已结业 / 3 已退课
	Status uint8

	EnrolledAt  time.Time    // datetime    报名时间
	CompletedAt sql.NullTime // datetime    结业时间（未结业为 NULL）
	DroppedAt   sql.NullTime // datetime    退课时间（未退课为 NULL）
	UpdatedAt   time.Time    // datetime    更新时间
	LockVersion uint64       // bigint unsigned 乐观锁版本号
}
