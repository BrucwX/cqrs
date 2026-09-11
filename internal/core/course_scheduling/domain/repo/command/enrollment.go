package command

import (
	"context"
	"errors"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
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
	Save(e *enrollment.CourseEnrollment) error
	// Delete 删除课程注册
	Delete(id int64) error
	// Enroll 学员选课
	//
	// checkConflictFn 由调用方注入，仓库会把「待写入的注册记录 + 该课程聚合」
	// 传进去；传 nil 表示不做检查。冲突时不写入。
	Enroll(ctx context.Context, e *enrollment.CourseEnrollment, checkConflictFn func(ctx context.Context, e *enrollment.CourseEnrollment, c course.Course) (bool, error)) error
}
