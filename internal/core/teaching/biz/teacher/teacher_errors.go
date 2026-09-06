package teacher

import "github.com/go-kratos/kratos/v3/errors"

var (
	// ErrTeacherNotFound is returned when a teacher does not exist.
	ErrTeacherNotFound = errors.NotFound("TEACHER_NOT_FOUND", "teacher not found")
	// ErrTeacherInvalidArgument is returned when a teacher request is invalid.
	ErrTeacherInvalidArgument = errors.BadRequest("TEACHER_INVALID_ARGUMENT", "invalid teacher argument")
)
