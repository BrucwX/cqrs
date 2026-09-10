package courseSlotChange

import "time"

// OriginalPlan 原始计划快照（用于溯源核对）
type OriginalPlan struct {
	slotID       int64     // 引用的课表模板 ID
	date         time.Time // 原定上课日期
	teacherID    int64     // 原讲师 ID
	classroomID  string    // 原教室 ID
	startTimeStr string    // 原起始时间（如 "16:00"）
	endTimeStr   string    // 原结束时间（如 "18:00"）
}

func NewOriginalPlan(
	slotID int64,
	date time.Time,
	teacherID int64,
	classroomID string,
	startTimeStr, endTimeStr string,
) OriginalPlan {
	return OriginalPlan{
		slotID:       slotID,
		date:         date,
		teacherID:    teacherID,
		classroomID:  classroomID,
		startTimeStr: startTimeStr,
		endTimeStr:   endTimeStr,
	}
}

func (op OriginalPlan) SlotID() int64        { return op.slotID }
func (op OriginalPlan) Date() time.Time      { return op.date }
func (op OriginalPlan) TeacherID() int64     { return op.teacherID }
func (op OriginalPlan) ClassroomID() string  { return op.classroomID }
func (op OriginalPlan) StartTimeStr() string { return op.startTimeStr }
func (op OriginalPlan) EndTimeStr() string   { return op.endTimeStr }
