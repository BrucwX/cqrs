package course

import "time"

// EnrollmentWindow 选课窗口与退课时效值对象
type EnrollmentWindow struct {
	startAt      time.Time
	endAt        time.Time
	dropDeadline time.Time
}

func NewEnrollmentWindow(start, end, drop time.Time) EnrollmentWindow {
	return EnrollmentWindow{startAt: start, endAt: end, dropDeadline: drop}
}

func (w EnrollmentWindow) CanEnroll(now time.Time) bool {
	return (now.Equal(w.startAt) || now.After(w.startAt)) && now.Before(w.endAt)
}

func (w EnrollmentWindow) CanDrop(now time.Time) bool {
	return now.Before(w.dropDeadline) || now.Equal(w.dropDeadline)
}
