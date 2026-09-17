package teacher

import "errors"

var (
	ErrInvalidPhone      = errors.New("invalid teacher phone format")
	ErrInvalidEmail      = errors.New("invalid teacher email format")
	ErrTeacherNotActive  = errors.New("teacher is not active (resigned or terminated)")
	ErrTeacherOnLeave    = errors.New("teacher is currently on leave")
	ErrAlreadyOnLeave    = errors.New("teacher is already on leave")
	ErrAlreadyActive     = errors.New("teacher is already in active status")
	ErrAlreadyTerminated = errors.New("teacher has already been terminated")
	ErrUnknownStatus     = errors.New("unknown teacher status")

	// ErrTeacherRequired 传入的讲师为空。
	ErrTeacherRequired = errors.New("teacher is required")
	// ErrTeacherNotFound 指定的讲师不存在。
	ErrTeacherNotFound = errors.New("teacher not found")
)
