package qualification

import "errors"

var (
	ErrQualificationExpired = errors.New("teacher qualification for this course has expired")
	ErrQualificationRevoked = errors.New("teacher qualification for this course has been revoked")
)
