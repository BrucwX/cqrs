package enrollment

import "errors"

var (
	ErrEnrollmentNotActive = errors.New("enrollment is not in active status")
	ErrCannotDropCompleted = errors.New("cannot drop a completed course")
	ErrAlreadyDropped      = errors.New("course enrollment has already been dropped")
	ErrAlreadyCompleted    = errors.New("course enrollment has already been completed")

	// ErrEnrollmentRequired 传入的注册记录为空。
	ErrEnrollmentRequired = errors.New("course enrollment is required")
	// ErrEnrollmentNotFound 指定的注册记录不存在。
	ErrEnrollmentNotFound = errors.New("course enrollment not found")
	// ErrEnrollmentConflict 选课被拒绝（不在选课窗口 / 已满 / 与在学课程撞时间）。
	ErrEnrollmentConflict = errors.New("course enrollment conflicts with existing constraints")
)
