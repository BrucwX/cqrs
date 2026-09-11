package courseSlotChange

import "errors"

var (
	ErrInvalidTimeRange = errors.New("target start time must be before target end time")
	ErrTargetCrossDay   = errors.New("target schedule must start and end on the same day")
	ErrTargetDateInPast = errors.New("target schedule date cannot be in the past")
)
