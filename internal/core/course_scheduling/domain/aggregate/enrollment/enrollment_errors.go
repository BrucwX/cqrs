package enrollment

import "errors"

var (
	ErrEnrollmentNotActive = errors.New("enrollment is not in active status")
	ErrCannotDropCompleted = errors.New("cannot drop a completed course")
	ErrAlreadyDropped      = errors.New("course enrollment has already been dropped")
	ErrAlreadyCompleted    = errors.New("course enrollment has already been completed")

	// ErrEnrollmentNotSelectable 该记录当前状态不能选课（只有「未选课」的记录能选）。
	ErrEnrollmentNotSelectable = errors.New("enrollment is not awaiting course selection")

	// ErrEnrollmentRequired 传入的注册记录为空。
	ErrEnrollmentRequired = errors.New("course enrollment is required")
	// ErrEnrollmentNotFound 指定的注册记录不存在。
	ErrEnrollmentNotFound = errors.New("course enrollment not found")
	// ErrEnrollmentNotPaid 没有该学员在该课程下「未选课」的注册记录
	// ⇒ 这门课还没缴费（没报班），不能选。
	ErrEnrollmentNotPaid = errors.New("student has not paid for this course")
	// ErrEnrollmentConflict 选课被拒绝（不在选课窗口 / 已满 / 与在学课程撞时间）。
	ErrEnrollmentConflict = errors.New("course enrollment conflicts with existing constraints")
)
