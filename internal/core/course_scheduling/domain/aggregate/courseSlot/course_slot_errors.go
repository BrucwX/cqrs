package courseSlot

import "errors"

var (
	ErrInvalidDayTime       = errors.New("start time of day must be before end time of day")
	ErrInvalidDurationHours = errors.New("duration must be positive")
	ErrSlotTimeOverlap      = errors.New("course slot template has time overlap on the same weekday")

	// ErrCourseSlotRequired 传入的课表槽位为空。
	ErrCourseSlotRequired = errors.New("course slot is required")
	// ErrCourseSlotNotFound 指定的课表槽位不存在。
	ErrCourseSlotNotFound = errors.New("course slot not found")
	// ErrCourseSlotConflict 目标讲师/教室/课程在该槽位的时间上已有其他安排。
	ErrCourseSlotConflict = errors.New("course slot conflicts with an existing schedule")
)
