package model

import (
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
)

// CourseSlot 对应表 course_slot，领域模型是 courseSlot.CourseSlot。
//
// 值对象 DayTimeRange 就地打平为 StartTime / EndTime 两个 MySQL TIME 列。
// TIME 列驱动以文本返回，所以这里用 string 承载，格式 "15:04:05"（领域模型里
// 是 DayTime：一天内的第几分钟）。
//
// 两个哨兵值表达「待分配」，与领域模型里的常量一致：
//
//	TeacherID   = -1  待分配讲师（CourseSlot 里是 PendingTeacherID）
//	ClassroomID = ""  待分配教室（CourseSlot 里是 PendingClassroomID）
//
// 正因为是哨兵值，这两列不适合加外键约束。
type CourseSlot struct {
	ID       string // varchar(36) 槽位 ID（uuid）
	CourseID string // varchar(36) 所属课程（course.id）

	// uint8 tinyint unsigned 星期几（0=周日 .. 6=周六，对应 time.Weekday）
	Weekday uint8

	StartTime string // time 当天上课开始时间（"15:04:05"）<- DayTimeRange.Start
	EndTime   string // time 当天上课结束时间（"15:04:05"）<- DayTimeRange.End

	TeacherID   int64  // bigint      默认讲师 ID（-1 = 待分配）
	ClassroomID string // varchar(36) 默认教室 ID（空 = 待分配）

	CreatedAt   time.Time // datetime 创建时间
	UpdatedAt   time.Time // datetime 更新时间
	LockVersion uint64    // bigint unsigned 乐观锁版本号
}

// CourseSlotToPO 写路径：DayTimeRange 打平成两个 time 列。
func CourseSlotToPO(do *courseSlot.CourseSlot) *CourseSlot {
	timeRange := do.TimeRange()
	return &CourseSlot{
		ID:          do.ID(),
		CourseID:    do.CourseID(),
		Weekday:     uint8(do.Weekday()),
		StartTime:   dayTimeToPO(timeRange.Start()),
		EndTime:     dayTimeToPO(timeRange.End()),
		TeacherID:   do.TeacherID(),
		ClassroomID: do.ClassroomID(),
		CreatedAt:   do.CreatedAt(),
		UpdatedAt:   do.UpdatedAt(),
	}
}

// CourseSlotToDO 读路径：两个 time 列拼回 DayTimeRange，先后顺序反了会被拒绝。
func CourseSlotToDO(po *CourseSlot) (*courseSlot.CourseSlot, error) {
	start, err := dayTimeFromPO(po.StartTime)
	if err != nil {
		return nil, err
	}
	end, err := dayTimeFromPO(po.EndTime)
	if err != nil {
		return nil, err
	}
	timeRange, err := courseSlot.NewDayTimeRange(start, end)
	if err != nil {
		return nil, err
	}
	return courseSlot.Reconstitute(
		po.ID, po.CourseID, time.Weekday(po.Weekday), timeRange,
		po.TeacherID, po.ClassroomID, po.CreatedAt, po.UpdatedAt,
	), nil
}
