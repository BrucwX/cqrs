package payment

import "github.com/go-kratos/kratos/v3/errors"

var (
	// ErrPaymentNotFound is returned when a payment does not exist.
	ErrPaymentNotFound = errors.NotFound("PAYMENT_NOT_FOUND", "payment not found")
	// ErrPaymentInvalidArgument is returned when a payment request is invalid.
	ErrPaymentInvalidArgument = errors.BadRequest("PAYMENT_INVALID_ARGUMENT", "invalid payment argument")
	// ErrPaymentUserNotFound is returned when the paying user does not exist.
	ErrPaymentUserNotFound = errors.BadRequest("PAYMENT_USER_NOT_FOUND", "user not found for payment")
	// ErrPaymentCourseNotFound is returned when the course to pay for does not exist.
	ErrPaymentCourseNotFound = errors.BadRequest("PAYMENT_COURSE_NOT_FOUND", "course not found for payment")
	// ErrPaymentInvalidStatus is returned when a payment status transition is invalid.
	ErrPaymentInvalidStatus = errors.BadRequest("PAYMENT_INVALID_STATUS", "invalid payment status for this operation")
)
