package course

import (
	"time"
)

// Semester is a value object representing a semester period with start and end dates.
type Semester struct {
	StartTime time.Time // 学期开始日期
	EndTime   time.Time // 学期结束日期
}

// NewSemester creates a new Semester value object.
// It validates that both times are valid and end is after start.
func NewSemester(start, end time.Time) (*Semester, error) {
	if start.IsZero() || end.IsZero() {
		return nil, ErrCourseInvalidArgument
	}
	if !end.After(start) {
		return nil, ErrCourseInvalidArgument
	}
	return &Semester{
		StartTime: start,
		EndTime:   end,
	}, nil
}

// Contains checks if the given time falls within the semester period.
func (s *Semester) Contains(t time.Time) bool {
	return !t.Before(s.StartTime) && !t.After(s.EndTime)
}

// Overlaps checks if two semesters have overlapping date ranges.
func (s *Semester) Overlaps(other *Semester) bool {
	if other == nil {
		return false
	}
	return !s.EndTime.Before(other.StartTime) && !other.EndTime.Before(s.StartTime)
}

// Duration returns the duration of the semester.
func (s *Semester) Duration() time.Duration {
	return s.EndTime.Sub(s.StartTime)
}
