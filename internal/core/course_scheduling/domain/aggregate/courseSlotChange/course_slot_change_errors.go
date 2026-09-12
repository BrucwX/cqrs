package courseSlotChange

import "errors"

var (
	ErrInvalidTimeRange = errors.New("target start time must be before target end time")
	ErrTargetCrossDay   = errors.New("target schedule must start and end on the same day")
	ErrTargetDateInPast = errors.New("target schedule date cannot be in the past")

	// ErrSlotChangeRequired 传入的课表变更单为空。
	ErrSlotChangeRequired = errors.New("course slot change is required")
	// ErrSlotChangeNotFound 指定的课表变更单不存在。
	ErrSlotChangeNotFound = errors.New("course slot change not found")
	// ErrSlotChangeConflict 换课被拒绝（目标讲师或教室在目标时段已被占用）。
	ErrSlotChangeConflict = errors.New("course slot change conflicts with existing schedule")
)
