package courseSlotChange

import "errors"

var (
	ErrInvalidTimeRange      = errors.New("target start time must be before target end time")
	ErrChangeAlreadyReviewed = errors.New("the slot change request has already been reviewed")
	ErrCannotCancelReviewed  = errors.New("cannot cancel a request that has already been approved or rejected")
	ErrTargetDateInPast      = errors.New("target schedule date cannot be in the past")
)
