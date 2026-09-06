package student

import "github.com/go-kratos/kratos/v3/errors"

var (
	// ErrStudentNotFound is returned when a student does not exist.
	ErrStudentNotFound = errors.NotFound("STUDENT_NOT_FOUND", "student not found")
	// ErrStudentInvalidArgument is returned when a student request is invalid.
	ErrStudentInvalidArgument = errors.BadRequest("STUDENT_INVALID_ARGUMENT", "invalid student argument")
)
