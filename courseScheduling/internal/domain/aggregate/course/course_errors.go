package course

import "errors"

var (
	ErrCourseFull           = errors.New("course capacity is full")
	ErrNotInEnrollmentStage = errors.New("not within enrollment time window")
	ErrDropDeadlinePassed   = errors.New("course drop deadline has passed")
	ErrInvalidHours         = errors.New("completed hours exceed total hours")
	ErrCourseTypeRequired   = errors.New("course type ID is required")
	ErrCapacityRequired     = errors.New("course capacity is required")
	ErrEnrollmentRequired   = errors.New("course enrollment window is required")
	ErrPeriodRequired       = errors.New("course period is required")

	// ErrCourseRequired 传入的课程为空。
	ErrCourseRequired = errors.New("course is required")
	// ErrCourseNotFound 指定的课程不存在。
	ErrCourseNotFound = errors.New("course not found")
)
