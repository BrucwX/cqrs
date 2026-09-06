package enrollment

import (
	"github.com/go-kratos/kratos/v3/errors"
)

var (
	ErrEnrollmentInvalidArgument = errors.BadRequest("ENROLLMENT_INVALID_ARGUMENT", "invalid enrollment argument")
	ErrEnrollmentNotFound        = errors.NotFound("ENROLLMENT_NOT_FOUND", "enrollment not found")
	ErrEnrollmentAlreadyActive   = errors.BadRequest("ENROLLMENT_ALREADY_ACTIVE", "enrollment already active")
	ErrEnrollmentCanceled        = errors.BadRequest("ENROLLMENT_CANCELED", "enrollment is canceled")
)
