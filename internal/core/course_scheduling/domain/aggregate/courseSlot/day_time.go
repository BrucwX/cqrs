package courseSlot

import (
	"errors"
	"fmt"
)

// DayTime 一天内的某个具体时间点（时、分）
type DayTime struct {
	hour   int // 0-23
	minute int // 0-59
}

func NewDayTime(hour, minute int) (DayTime, error) {
	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return DayTime{}, errors.New("invalid hour or minute range")
	}
	return DayTime{hour: hour, minute: minute}, nil
}

func (dt DayTime) TotalMinutes() int {
	return dt.hour*60 + dt.minute
}

func (dt DayTime) String() string {
	return fmt.Sprintf("%02d:%02d", dt.hour, dt.minute)
}
