package course

import "errors"

var (
	ErrCourseFull           = errors.New("course capacity is full")
	ErrNotInEnrollmentStage = errors.New("not within enrollment time window")
	ErrDropDeadlinePassed   = errors.New("course drop deadline has passed")
	ErrInvalidHours         = errors.New("completed hours exceed total hours")
)
