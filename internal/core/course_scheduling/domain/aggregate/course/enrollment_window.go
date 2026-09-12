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

// --- 只读属性访问器 (Getters) ---

func (w EnrollmentWindow) StartAt() time.Time      { return w.startAt }
func (w EnrollmentWindow) EndAt() time.Time        { return w.endAt }
func (w EnrollmentWindow) DropDeadline() time.Time { return w.dropDeadline }
