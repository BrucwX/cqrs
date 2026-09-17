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
	// 只管写：把这条注册记录落库（未选课 -> 在读 的状态流转由聚合负责，
	// 资格与准入判定由调用方在调过来之前做完 —— 见
	// domain/service/scheduleConflict 的 Service.CheckEnrollment）。
	// 所以这里没有回调，也没有「传 nil 表示不检查」这类分支。
	Enroll(ctx context.Context, e *enrollment.CourseEnrollment) error
	// GetUnSelectEnrollBySC 取该学员（studentID）在该课程（courseID）下「未选课」的
	// 那条注册记录 —— 也就是他的选课资格。
	//
	// 同一门课可能有多条记录（退课后重新登记会再多一条），取 ID 最小的那条。
	// 取不到时报 ErrEnrollmentNotPaid —— 没有这条记录就是这门课还没缴费，
	// 所以直接把「未付费」报出来，调用方不用再做二次翻译。
	GetUnSelectEnrollBySC(ctx context.Context, studentID int64, courseID string) (enrollment.CourseEnrollment, error)
	// GetEnrollments 取该学员的全部报名记录
	GetEnrollments(ctx context.Context, studentID int64) ([]enrollment.CourseEnrollment, error)
}
