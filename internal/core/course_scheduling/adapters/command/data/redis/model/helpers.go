package model

import (
	"fmt"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
)

// 本文件放 model 包内部共用的转换零件。
//
// 每个聚合的 XxxToRedis / XxxFromRedis 都写在它自己的文件里
// （course_slot.go 里就是 CourseSlotToRedis / CourseSlotFromRedis）。
// 这里只留被多个文件复用的部分：
//
//	timePtr / timeValue   json null <-> 零值时间
//	dayTimeToRedis / dayTimeFromRedis  "HH:MM" <-> courseSlot.DayTime 值对象

// timePtr 把零值时间写成 JSON null —— 领域模型用零值表示「没有」。
func timePtr(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

// timeValue 把 JSON null 读回零值时间。
func timeValue(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

// dayTimeToRedis 把 DayTime 值对象写成 "HH:MM"。
func dayTimeToRedis(t courseSlot.DayTime) string {
	minutes := t.TotalMinutes()
	return fmt.Sprintf("%02d:%02d", minutes/60, minutes%60)
}

// dayTimeFromRedis 解析 "HH:MM"。
func dayTimeFromRedis(s string) (courseSlot.DayTime, error) {
	parsed, err := time.Parse("15:04", s)
	if err != nil {
		return courseSlot.DayTime{}, fmt.Errorf("model: 无法解析时间 %q", s)
	}
	return courseSlot.NewDayTime(parsed.Hour(), parsed.Minute())
}
