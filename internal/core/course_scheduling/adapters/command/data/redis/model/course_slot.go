package model

import (
	"fmt"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
)

// CourseSlot 是 Redis 快照（对应 course_slot），领域模型是 courseSlot.CourseSlot。
//
// 值对象 DayTimeRange 就地打平为 StartTime / EndTime 两个字符串。
// 这里用 string 承载，格式 "15:04:05"（领域模型里
// 是 DayTime：一天内的第几分钟）。
//
// 两个哨兵值表达「待分配」，与领域模型里的常量一致：
//
//	TeacherID   = -1  待分配讲师（CourseSlot 里是 PendingTeacherID）
//	ClassroomID = ""  待分配教室（CourseSlot 里是 PendingClassroomID）
//
// 正因为是哨兵值，这两个字段没有外键约束。
type CourseSlot struct {
	ID       string `json:"id"`        // 槽位 ID（uuid）
	CourseID string `json:"course_id"` // 所属课程（course.id）

	// 星期几（0=周日 .. 6=周六，对应 time.Weekday）
	Weekday uint8 `json:"weekday"`

	StartTime string `json:"start_time"` // 当天上课开始时间（"15:04:05"）<- DayTimeRange.Start
	EndTime   string `json:"end_time"`   // 当天上课结束时间（"15:04:05"）<- DayTimeRange.End

	TeacherID   int64  `json:"teacher_id"`   // 默认讲师 ID（-1 = 待分配）
	ClassroomID string `json:"classroom_id"` // 默认教室 ID（空 = 待分配）

	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// CourseSlotToRedis 写路径：DayTimeRange 打平成两个时间字符串。
func CourseSlotToRedis(do *courseSlot.CourseSlot) (*CourseSlot, error) {
	if do == nil {
		return nil, ErrCourseSlotToRedis
	}
	timeRange := do.TimeRange()
	return &CourseSlot{
		ID:          do.ID(),
		CourseID:    do.CourseID(),
		Weekday:     uint8(do.Weekday()),
		StartTime:   dayTimeToRedis(timeRange.Start()),
		EndTime:     dayTimeToRedis(timeRange.End()),
		TeacherID:   do.TeacherID(),
		ClassroomID: do.ClassroomID(),
		CreatedAt:   do.CreatedAt(),
		UpdatedAt:   do.UpdatedAt(),
	}, nil
}

// CourseSlotFromRedis 读路径：两个时间字符串拼回 DayTimeRange，先后顺序反了会被拒绝。
func CourseSlotFromRedis(po *CourseSlot) (*courseSlot.CourseSlot, error) {
	if po == nil {
		return nil, ErrCourseSlotFromRedis
	}
	start, err := dayTimeFromRedis(po.StartTime)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCourseSlotFromRedis, err)
	}
	end, err := dayTimeFromRedis(po.EndTime)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCourseSlotFromRedis, err)
	}
	timeRange, err := courseSlot.NewDayTimeRange(start, end)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCourseSlotFromRedis, err)
	}
	return courseSlot.Reconstitute(
		po.ID, po.CourseID, time.Weekday(po.Weekday), timeRange,
		po.TeacherID, po.ClassroomID, po.CreatedAt, po.UpdatedAt,
	), nil
}
