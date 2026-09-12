package command

import (
	"context"
	"errors"

	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
)

var (
	// ErrEnrollmentRequired 传入的注册记录为空。
	ErrEnrollmentRequired = errors.New("course enrollment is required")
	// ErrEnrollmentNotFound 指定的注册记录不存在。
	ErrEnrollmentNotFound = errors.New("course enrollment not found")
	// ErrEnrollmentConflict 选课被拒绝（不在选课窗口 / 已满 / 与在学课程撞时间）。
	ErrEnrollmentConflict = errors.New("course enrollment conflicts with existing constraints")
)

// CourseEnrollmentCommand 课程注册命令接口
type CourseEnrollmentCommand interface {
	// Save 保存课程注册（新增或更新，不做检查）
	Save(ctx context.Context, e *enrollment.CourseEnrollment) error
	// Delete 删除课程注册
	Delete(ctx context.Context, id int64) error
	// MustGet 取课程注册聚合；不存在时报 ErrEnrollmentNotFound
	MustGet(ctx context.Context, id int64) (enrollment.CourseEnrollment, error)
	// Enroll 学员选课
	//
	// 只管写：把这条注册记录落库。准入判定（选课窗口 / 容量 / 时间冲突）
	// 不在这里 —— 那是调用方的事（见 domain/service/schedule 的 Conflict.CheckEnrollment），
	// 判定通过才调过来。所以这里没有回调，也没有「传 nil 表示不检查」这类分支。
	Enroll(ctx context.Context, e *enrollment.CourseEnrollment) error
	// GetEnrollments 取该学员的全部报名记录
	GetEnrollments(ctx context.Context, studentID int64) ([]enrollment.CourseEnrollment, error)
}
