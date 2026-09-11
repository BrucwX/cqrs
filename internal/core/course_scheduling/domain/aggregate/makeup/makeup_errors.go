package makeup

import "errors"

var (
	ErrCourseMismatch   = errors.New("target course must be identical to the missed course")
	ErrAlreadyFinalized = errors.New("makeup is already completed or cancelled")
)
