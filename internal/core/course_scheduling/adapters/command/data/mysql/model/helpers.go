package model

import (
	"database/sql"
	"fmt"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
)

// 本文件放 model 包内部共用的转换零件：可空时间列与 DayTime 列的互转。
//
// 每个聚合的 XxxToPO / XxxToDO 都写在它自己那张表的文件里
// （course_slot.go 里就是 CourseSlotToPO / CourseSlotToDO）。
// 这里只留被多个表文件复用的部分：
//
//	nullTime / timeValue          datetime 可空列 <-> 零值时间
//	dayTimeToPO / dayTimeFromPO   time 列 <-> courseSlot.DayTime 值对象

// nullTime 把零值时间写成 NULL —— 领域模型用零值表示「没有」，MySQL 用 NULL。
func nullTime(t time.Time) sql.NullTime {
	return sql.NullTime{Time: t, Valid: !t.IsZero()}
}

// timeValue 把 NULL 读回零值时间。
func timeValue(n sql.NullTime) time.Time {
	if !n.Valid {
		return time.Time{}
	}
	return n.Time
}

// dayTimeToPO 把 DayTime 值对象写成 "15:04:05"。
func dayTimeToPO(t courseSlot.DayTime) string {
	minutes := t.TotalMinutes()
	return fmt.Sprintf("%02d:%02d:00", minutes/60, minutes%60)
}

// dayTimeFromPO 解析 time 列。驱动可能回 "15:04:05" 也可能回 "15:04"
// （取决于列类型与 parseTime 设置），两种都收。
func dayTimeFromPO(s string) (courseSlot.DayTime, error) {
	for _, layout := range []string{"15:04:05", "15:04"} {
		parsed, err := time.Parse(layout, s)
		if err == nil {
			return courseSlot.NewDayTime(parsed.Hour(), parsed.Minute())
		}
	}
	return courseSlot.DayTime{}, fmt.Errorf("model: 无法解析时间列 %q", s)
}
