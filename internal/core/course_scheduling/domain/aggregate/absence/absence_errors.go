package absence

import "errors"

var (
	ErrInvalidHours = errors.New("missed hours must be greater than zero")
)
