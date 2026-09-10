package absence

import "errors"

var (
	ErrInvalidHours           = errors.New("missed hours must be greater than zero")
	ErrAlreadyReviewed        = errors.New("absence record has already been reviewed")
	ErrAlreadyRectified       = errors.New("absence has already been rectified or made up")
	ErrCannotRectifyUnexcused = errors.New("unexcused absence cannot be rectified without appeal")
)
