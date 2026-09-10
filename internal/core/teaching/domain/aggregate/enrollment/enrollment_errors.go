package enrollment

import "errors"

var (
	ErrEnrollmentNotActive = errors.New("enrollment is not in active status")
	ErrCannotDropCompleted = errors.New("cannot drop a completed course")
	ErrAlreadyDropped      = errors.New("course enrollment has already been dropped")
	ErrAlreadyCompleted    = errors.New("course enrollment has already been completed")
)
