package makeup

import "errors"

var (
	ErrAlreadyFinalized = errors.New("makeup is already completed or cancelled")
)
