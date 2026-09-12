package classroom

import "errors"

var (
	ErrInvalidCapacity         = errors.New("classroom capacity must be greater than zero")
	ErrInvalidLocation         = errors.New("classroom location is required")
	ErrClassroomUnavailable    = errors.New("classroom is currently under maintenance or disabled")
	ErrCapacityExceeded        = errors.New("required seats exceed classroom capacity")
	ErrUnsupportedStatusChange = errors.New("classroom status can only be switched between available and under maintenance")

	// ErrClassroomRequired 传入的教室为空。
	ErrClassroomRequired = errors.New("classroom is required")
	// ErrClassroomNotFound 指定的教室不存在。
	ErrClassroomNotFound = errors.New("classroom not found")
)
