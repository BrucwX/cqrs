package courseSlot

// DayTimeRange 一天内的上课时间段（如 16:00 - 18:00）
type DayTimeRange struct {
	start DayTime
	end   DayTime
}

func NewDayTimeRange(start, end DayTime) (DayTimeRange, error) {
	if start.TotalMinutes() >= end.TotalMinutes() {
		return DayTimeRange{}, ErrInvalidDayTime
	}
	return DayTimeRange{start: start, end: end}, nil
}

func (dtr DayTimeRange) Start() DayTime { return dtr.start }
func (dtr DayTimeRange) End() DayTime   { return dtr.end }

// DurationHours 本时间段折算的课时量（如 2.0 小时）
func (dtr DayTimeRange) DurationHours() float64 {
	return float64(dtr.end.TotalMinutes()-dtr.start.TotalMinutes()) / 60.0
}

// Overlaps 校验同一天内的两个时间区间是否有重叠
func (dtr DayTimeRange) Overlaps(other DayTimeRange) bool {
	return dtr.start.TotalMinutes() < other.end.TotalMinutes() &&
		other.start.TotalMinutes() < dtr.end.TotalMinutes()
}
