package makeup

import "errors"

var (
	ErrAlreadyFinalized = errors.New("makeup is already completed or cancelled")

	// ErrMakeupRequired 传入的补课预约为空。
	ErrMakeupRequired = errors.New("student makeup is required")
	// ErrMakeupNotFound 指定的补课预约不存在。
	ErrMakeupNotFound = errors.New("student makeup not found")
	// ErrMakeupConflict 补课预约被拒绝（目标那节课的教室装不下）。
	ErrMakeupConflict = errors.New("student makeup conflicts with classroom capacity")
)
