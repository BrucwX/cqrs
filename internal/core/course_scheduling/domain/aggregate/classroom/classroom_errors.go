package classroom

import "errors"

var (
	ErrInvalidCapacity      = errors.New("classroom capacity must be greater than zero")
	ErrClassroomUnavailable = errors.New("classroom is currently under maintenance or disabled")
	ErrCapacityExceeded     = errors.New("required seats exceed classroom capacity")
)
