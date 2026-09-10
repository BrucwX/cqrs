package course

import (
	"errors"
	"time"
)

// CoursePeriod 教学周期与课时值对象
type CoursePeriod struct {
	startAt        time.Time
	endAt          time.Time
	totalHours     int
	completedHours int
}

func NewCoursePeriod(start, end time.Time, total int, completed int) (CoursePeriod, error) {
	if completed > total || total <= 0 || completed < 0 {
		return CoursePeriod{}, errors.New("completed hours exceed total hours")
	}
	return CoursePeriod{
		startAt:        start,
		endAt:          end,
		totalHours:     total,
		completedHours: completed,
	}, nil
}
