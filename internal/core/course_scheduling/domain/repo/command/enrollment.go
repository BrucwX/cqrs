package command

import (
	"context"
	"errors"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
)

var (
	// ErrEnrollmentRequired 传入的注册记录为空。
	ErrEnrollmentRequired = errors.New("course enrollment is required")
	// ErrEnrollmentNotFound 指定的注册记录不存在。
	ErrEnrollmentNotFound = errors.New("course enrollment not found")
	// ErrEnrollmentConflict 选课被拒绝（不在选课窗口 / 已满 / 与在学课程撞时间，
	// 或已选过该课且该课有排期）。
	ErrEnrollmentConflict = errors.New("course enrollment conflicts with existing constraints")
)

// EnrollmentContext 是仓库侧在选课时一并交给冲突检查的上下文。
//
// 仓库（命令适配器）能直接读到读模型，所以由它把 checkEnrollment 需要的
// 「目标课程 / 目标课程排期 / 学员现有排期」一次装好传进来。
type EnrollmentContext struct {
	// Course 目标课程（用于判定选课窗口与容量）；取不到时为零值
	Course course.Course
	// TargetSlots 目标课程的排期
	TargetSlots courseSlot.CourseSlots
	// StudentSlots 该学员现有在学课程的排期（含目标课程本身，
	// 所以重复选课会因自身槽位重叠而被判为时间冲突）
	StudentSlots courseSlot.CourseSlots
}

// CourseEnrollmentCommand 课程注册命令接口
type CourseEnrollmentCommand interface {
	// Save 保存课程注册（新增或更新，不做检查）
	Save(e *enrollment.CourseEnrollment) error
	// Delete 删除课程注册
	Delete(id int64) error
	// Enroll 学员选课
	//
	// checkConflictFn 由调用方注入，仓库会把 EnrollmentContext 装好传进去；
	// 传 nil 表示不做检查。冲突时不写入。
	Enroll(ctx context.Context, e *enrollment.CourseEnrollment, checkConflictFn func(ctx context.Context, ec EnrollmentContext) (bool, error)) error
}
