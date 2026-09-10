package makeup

import "errors"

var (
	ErrCourseMismatch       = errors.New("target course must be identical to the missed course")
	ErrAlreadyReviewed      = errors.New("makeup request has already been reviewed")
	ErrCannotCompleteNotApp = errors.New("cannot complete attendance on an unapproved makeup request")
	ErrAlreadyFinalized     = errors.New("makeup request is already completed or cancelled")
)
