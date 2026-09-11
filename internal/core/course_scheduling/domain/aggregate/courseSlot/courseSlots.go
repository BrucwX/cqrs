package courseSlot

import "time"

// CourseSlots 是一组周排期槽位，代表某个讲师 / 教室 / 课程占用的时间集合。
type CourseSlots []CourseSlot

// ConflictsWith 判断两组槽位之间是否存在时间冲突。
//
// 任一对槽位「同星期几 + 时间段重叠」即视为冲突；
// 判定只看时间，不看讲师/教室（那属于 IsConflictingWith 的语义）。
func (s CourseSlots) ConflictsWith(other CourseSlots) bool {
	for _, a := range s {
		for _, b := range other {
			if a.OverlapsWith(b) {
				return true
			}
		}
	}
	return false
}

// OverlapsOn 判断组内是否有槽位落在一个具体的「星期几 + 时间段」上。
//
// 用于拿具体某一天的上课时间（如某次临时换课的目标时段）去比对周排期：
// 同星期几 且 时间段重叠即算占用。
func (s CourseSlots) OverlapsOn(weekday time.Weekday, span DayTimeRange) bool {
	for _, cs := range s {
		if cs.weekday == weekday && cs.timeRange.Overlaps(span) {
			return true
		}
	}
	return false
}
