package courseSlot

import "errors"

var (
	ErrInvalidDayTime       = errors.New("start time of day must be before end time of day")
	ErrInvalidDurationHours = errors.New("duration must be positive")
	ErrSlotTimeOverlap      = errors.New("course slot template has time overlap on the same weekday")
)
