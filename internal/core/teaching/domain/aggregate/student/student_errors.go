package student

import "errors"

var (
	ErrInvalidPhone             = errors.New("invalid student phone format")
	ErrInvalidEmail             = errors.New("invalid student email format")
	ErrStudentBlocked           = errors.New("student account is blocked")
	ErrStudentCancelled         = errors.New("student account has been cancelled")
	ErrSelfEnrollmentProhibited = errors.New("a teacher cannot enroll in their own course as a student")
	ErrUnknownStatus            = errors.New("unknown student status")
)
