package courseSlotChange

import (
	"errors"
	"time"
)

// TargetPlan 调整后的实际执行计划
type TargetPlan struct {
	targetStartAt time.Time // 调整后具体开始时间（年月日 时分）
	targetEndAt   time.Time // 调整后具体结束时间（年月日 时分）
	teacherID     int64     // 实际授课讲师（可能是原讲师或代课老师）
	classroomID   string    // 实际使用教室
}

func NewTargetPlan(
	start, end time.Time,
	teacherID int64,
	classroomID string,
) (TargetPlan, error) {
	if !start.Before(end) {
		return TargetPlan{}, ErrInvalidTimeRange
	}
	if teacherID <= 0 {
		return TargetPlan{}, errors.New("target teacher ID is required")
	}
	if classroomID == "" {
		return TargetPlan{}, errors.New("target classroom ID is required")
	}

	return TargetPlan{
		targetStartAt: start,
		targetEndAt:   end,
		teacherID:     teacherID,
		classroomID:   classroomID,
	}, nil
}

func (tp TargetPlan) TargetStartAt() time.Time { return tp.targetStartAt }
func (tp TargetPlan) TargetEndAt() time.Time   { return tp.targetEndAt }
func (tp TargetPlan) TeacherID() int64         { return tp.teacherID }
func (tp TargetPlan) ClassroomID() string      { return tp.classroomID }
