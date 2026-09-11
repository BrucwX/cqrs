package courseSlotChange

import "errors"

var (
	ErrInvalidTimeRange = errors.New("target start time must be before target end time")
	ErrTargetDateInPast = errors.New("target schedule date cannot be in the past")
)
