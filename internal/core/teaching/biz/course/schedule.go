package course

import (
	"fmt"
	"time"
)

// Weekday represents a day of the week.
type Weekday int

const (
	WeekdayUnspecified Weekday = 0
	WeekdayMonday      Weekday = 1
	WeekdayTuesday     Weekday = 2
	WeekdayWednesday   Weekday = 3
	WeekdayThursday    Weekday = 4
	WeekdayFriday      Weekday = 5
	WeekdaySaturday    Weekday = 6
	WeekdaySunday      Weekday = 7
)

// Session represents a teaching block of the day (morning / afternoon / evening).
type Session int

const (
	SessionUnspecified Session = 0
	SessionMorning     Session = 1 // 上午
	SessionAfternoon   Session = 2 // 下午
	SessionEvening     Session = 3 // 晚上
)

// String returns the Chinese name of the session.
func (s Session) String() string {
	switch s {
	case SessionMorning:
		return "上午"
	case SessionAfternoon:
		return "下午"
	case SessionEvening:
		return "晚上"
	default:
		return "未知"
	}
}

// Period is a value object describing one fixed teaching period of the
// venue-wide standard timetable. Real times come from that table, so callers
// cannot create arbitrary time ranges — every course/classroom shares the same
// clock. Fields are exported so the ordinal, session and clock can be read
// directly.
type Period struct {
	// Number 是该时段内第几节课（从 1 起），如 上午第2节 → 2。
	Number int
	// Session 是该节次所属的时段（上午 / 下午 / 晚上）。
	Session Session
	// Start 是统一作息表规定的开始时间（只含时分，基日 0001-01-01）。
	Start time.Time
	// End 是统一作息表规定的结束时间（只含时分，基日 0001-01-01）。
	End time.Time
}

// newPeriod builds a canonical Period from clock strings.
func newPeriod(session Session, number int, start, end string) Period {
	sh, sm := splitClock(start)
	eh, em := splitClock(end)
	return Period{
		Number:  number,
		Session: session,
		Start:   time.Date(0, 1, 1, sh, sm, 0, 0, time.UTC),
		End:     time.Date(0, 1, 1, eh, em, 0, 0, time.UTC),
	}
}

// 标准节次（按全天上课顺序）。调整作息只需改这里。
var (
	PeriodAM1  = newPeriod(SessionMorning, 1, "08:00", "08:50")   // 上午第1节
	PeriodAM2  = newPeriod(SessionMorning, 2, "09:00", "09:50")   // 上午第2节
	PeriodAM3  = newPeriod(SessionMorning, 3, "10:00", "10:50")   // 上午第3节
	PeriodAM4  = newPeriod(SessionMorning, 4, "11:00", "11:50")   // 上午第4节
	PeriodPM1  = newPeriod(SessionAfternoon, 1, "14:00", "14:50") // 下午第1节
	PeriodPM2  = newPeriod(SessionAfternoon, 2, "15:00", "15:50") // 下午第2节
	PeriodPM3  = newPeriod(SessionAfternoon, 3, "16:00", "16:50") // 下午第3节
	PeriodEVE1 = newPeriod(SessionEvening, 1, "18:30", "19:20")   // 晚上第1节
	PeriodEVE2 = newPeriod(SessionEvening, 2, "19:30", "20:20")   // 晚上第2节
)

// standardTimetable 以全天上课顺序保存所有标准节次（唯一事实来源）。
var standardTimetable = []Period{
	PeriodAM1, PeriodAM2, PeriodAM3, PeriodAM4,
	PeriodPM1, PeriodPM2, PeriodPM3,
	PeriodEVE1, PeriodEVE2,
}

// Valid reports whether p is one of the standard periods.
func (p Period) Valid() bool {
	for _, s := range standardTimetable {
		if p == s {
			return true
		}
	}
	return false
}

// String returns a human-readable label, e.g. "上午第3节(10:00-10:50)".
func (p Period) String() string {
	return fmt.Sprintf("%s第%d节(%s-%s)",
		p.Session, p.Number, p.Start.Format("15:04"), p.End.Format("15:04"))
}

// AllPeriods returns every standard period in day order (morning -> evening).
func AllPeriods() []Period {
	return append([]Period(nil), standardTimetable...)
}

// splitClock parses an "HH:MM" clock into hour and minute.
func splitClock(clock string) (hour, min int) {
	_, _ = fmt.Sscanf(clock, "%d:%d", &hour, &min)
	return
}

// ScheduleTime is a closed value object representing a recurring teaching slot,
// e.g. Wednesday 上午第1节. Times come from the venue-wide standard timetable,
// so callers cannot create arbitrary time ranges.
type ScheduleTime struct {
	Day    Weekday // 星期几
	Period Period  // 节次（对应统一作息表）
}

// NewScheduleTime creates a ScheduleTime from a weekday and a fixed period.
func NewScheduleTime(day Weekday, period Period) (*ScheduleTime, error) {
	if day == WeekdayUnspecified {
		return nil, ErrCourseInvalidSchedule
	}
	if !period.Valid() {
		return nil, ErrCourseInvalidSchedule
	}
	return &ScheduleTime{Day: day, Period: period}, nil
}

// Session returns the teaching block of the slot.
func (s *ScheduleTime) Session() Session {
	return s.Period.Session
}

// StartTime returns the canonical start time (time-of-day only).
func (s *ScheduleTime) StartTime() time.Time {
	return s.Period.Start
}

// EndTime returns the canonical end time (time-of-day only).
func (s *ScheduleTime) EndTime() time.Time {
	return s.Period.End
}

// Duration returns the duration of the schedule time slot.
func (s *ScheduleTime) Duration() time.Duration {
	return s.EndTime().Sub(s.StartTime())
}

// String returns a human-readable representation, e.g. "周三 上午第1节(08:00-08:50)".
func (s *ScheduleTime) String() string {
	return fmt.Sprintf("%s %s", s.Day.String(), s.Period.String())
}

// Equals checks if two ScheduleTime value objects are equal.
func (s *ScheduleTime) Equals(other *ScheduleTime) bool {
	if other == nil {
		return false
	}
	return s.Day == other.Day && s.Period == other.Period
}

// Overlaps checks if two ScheduleTime slots would occupy the same time.
// Since all times come from the closed standard timetable, two slots overlap
// iff they fall on the same weekday and the same period.
func (s *ScheduleTime) Overlaps(other *ScheduleTime) bool {
	if other == nil {
		return false
	}
	return s.Day == other.Day && s.Period == other.Period
}

// String returns the Chinese name of the weekday.
func (w Weekday) String() string {
	switch w {
	case WeekdayMonday:
		return "周一"
	case WeekdayTuesday:
		return "周二"
	case WeekdayWednesday:
		return "周三"
	case WeekdayThursday:
		return "周四"
	case WeekdayFriday:
		return "周五"
	case WeekdaySaturday:
		return "周六"
	case WeekdaySunday:
		return "周日"
	default:
		return "未知"
	}
}
