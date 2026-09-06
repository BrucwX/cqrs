package classroom

import "github.com/go-kratos/kratos/v3/errors"

var (
	// ErrClassroomNotFound is returned when a classroom does not exist.
	ErrClassroomNotFound = errors.NotFound("CLASSROOM_NOT_FOUND", "classroom not found")
	// ErrClassroomInvalidArgument is returned when a classroom request is invalid.
	ErrClassroomInvalidArgument = errors.BadRequest("CLASSROOM_INVALID_ARGUMENT", "invalid classroom argument")
)
