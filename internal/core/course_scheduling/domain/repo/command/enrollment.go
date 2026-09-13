package command

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
)

// CourseEnrollmentCommand 课程注册命令接口
type CourseEnrollmentCommand interface {
	// SaveEnrollment 保存课程注册（新增或更新，不做检查）
	SaveEnrollment(ctx context.Context, e *enrollment.CourseEnrollment) error
	// DeleteEnrollment 删除课程注册
	DeleteEnrollment(ctx context.Context, id int64) error
	// MustGetEnrollment 取课程注册聚合；不存在时报 ErrEnrollmentNotFound
	MustGetEnrollment(ctx context.Context, id int64) (enrollment.CourseEnrollment, error)
	// Enroll 学员选课
	//
	// 只管写：把这条注册记录落库。准入判定（选课窗口 / 容量 / 时间冲突）
	// 不在这里 —— 那是调用方的事（见 domain/service/schedule 的 Conflict.CheckEnrollment），
	// 判定通过才调过来。所以这里没有回调，也没有「传 nil 表示不检查」这类分支。
	Enroll(ctx context.Context, e *enrollment.CourseEnrollment) error
	// GetEnrollments 取该学员的全部报名记录
	GetEnrollments(ctx context.Context, studentID int64) ([]enrollment.CourseEnrollment, error)
}
