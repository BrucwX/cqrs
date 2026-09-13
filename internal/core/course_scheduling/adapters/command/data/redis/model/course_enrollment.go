package model

import (
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
)

// CourseEnrollment 是 Redis 快照（对应 course_enrollment），领域模型是 enrollment.CourseEnrollment。
//
// 同一个学员可以重复报同一门课（退课后重选），所以**没有唯一键**，
// 履约过程由 Status 状态机表达（在读 / 已结业 / 已退课）。
//
// 两个可空时间在领域模型里是零值时间，在这里就是 null：
//
//	CompletedAt null = 未结业
//	DroppedAt   null = 未退课
type CourseEnrollment struct {
	ID        int64  `json:"id"`         // 报名记录 ID
	StudentID int64  `json:"student_id"` // 学员 ID
	CourseID  string `json:"course_id"`  // 课程 ID（course.id）

	// 1 在读 / 2 已结业 / 3 已退课
	Status uint8 `json:"status"`

	EnrolledAt  time.Time  `json:"enrolled_at"`  // 报名时间
	CompletedAt *time.Time `json:"completed_at"` // 结业时间（未结业为 null）
	DroppedAt   *time.Time `json:"dropped_at"`   // 退课时间（未退课为 null）
	UpdatedAt   time.Time  `json:"updated_at"`   // 更新时间
}

// EnrollmentToRedis 写路径：零值时间表示「没有」，写成 JSON 就是 null。
func EnrollmentToRedis(do *enrollment.CourseEnrollment) (*CourseEnrollment, error) {
	if do == nil {
		return nil, ErrEnrollmentToRedis
	}
	return &CourseEnrollment{
		ID:          do.ID(),
		StudentID:   do.StudentID(),
		CourseID:    do.CourseID(),
		Status:      uint8(do.Status()),
		EnrolledAt:  do.EnrolledAt(),
		CompletedAt: timePtr(do.CompletedAt()),
		DroppedAt:   timePtr(do.DroppedAt()),
		UpdatedAt:   do.UpdatedAt(),
	}, nil
}

// EnrollmentFromRedis 读路径：null 还原成零值时间。
func EnrollmentFromRedis(po *CourseEnrollment) (*enrollment.CourseEnrollment, error) {
	if po == nil {
		return nil, ErrEnrollmentFromRedis
	}
	return enrollment.Reconstitute(
		po.ID, po.StudentID, po.CourseID, enrollment.Status(po.Status),
		po.EnrolledAt, timeValue(po.CompletedAt), timeValue(po.DroppedAt), po.UpdatedAt,
	), nil
}
