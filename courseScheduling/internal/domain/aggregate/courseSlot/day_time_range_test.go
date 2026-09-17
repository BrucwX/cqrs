package courseSlot

import (
	"errors"
	"testing"
	"time"
)

// TestNewDayTime 一天内的时间点边界。
func TestNewDayTime(t *testing.T) {
	cases := []struct {
		name   string
		hour   int
		minute int
		wantOK bool
	}{
		{name: "00:00", hour: 0, minute: 0, wantOK: true},
		{name: "23:59", hour: 23, minute: 59, wantOK: true},
		{name: "12:30", hour: 12, minute: 30, wantOK: true},
		{name: "小时为负", hour: -1, minute: 0},
		{name: "小时越界", hour: 24, minute: 0},
		{name: "分钟为负", hour: 0, minute: -1},
		{name: "分钟越界", hour: 0, minute: 60},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewDayTime(tc.hour, tc.minute)
			if tc.wantOK {
				if err != nil {
					t.Fatalf("NewDayTime(%d, %d): %v", tc.hour, tc.minute, err)
				}
				if want := tc.hour*60 + tc.minute; got.TotalMinutes() != want {
					t.Errorf("TotalMinutes = %d, want %d", got.TotalMinutes(), want)
				}
				return
			}
			if err == nil {
				t.Errorf("NewDayTime(%d, %d) 应该被拒绝，实际得到 %v", tc.hour, tc.minute, got)
			}
		})
	}
}

// TestNewDayTimeRangeRejectsCrossDay 时间段不允许跨天。
//
// DayTimeRange 只有「时:分」两个端点，本身就没有「日期」的概念，
// 所以跨天在结构上就不可能出现；而任何「结束不晚于开始」的区间
// （含 23:00 -> 01:00 这种绕午夜的写法）都会被挡下来。
func TestNewDayTimeRangeRejectsCrossDay(t *testing.T) {
	at := func(t *testing.T, hour, minute int) DayTime {
		t.Helper()
		dt, err := NewDayTime(hour, minute)
		if err != nil {
			t.Fatalf("new day time: %v", err)
		}
		return dt
	}

	cases := []struct {
		name string
		from int
		to   int
	}{
		{name: "绕午夜 23:00 -> 01:00", from: 23, to: 1},
		{name: "绕午夜 22:00 -> 02:00", from: 22, to: 2},
		{name: "绕午夜 23:30 -> 00:30", from: 23, to: 0},
		{name: "起止相同", from: 9, to: 9},
		{name: "结束早于开始", from: 11, to: 9},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewDayTimeRange(at(t, tc.from, 0), at(t, tc.to, 0))
			if !errors.Is(err, ErrInvalidDayTime) {
				t.Errorf("err = %v, want %v", err, ErrInvalidDayTime)
			}
			if got != (DayTimeRange{}) {
				t.Errorf("被拒绝时应返回零值，实际 %v", got)
			}
		})
	}
}

// TestNewDayTimeRangeSameDay 当天内的合法区间都接受，最长为 00:00 -> 23:59。
func TestNewDayTimeRangeSameDay(t *testing.T) {
	at := func(t *testing.T, hour, minute int) DayTime {
		t.Helper()
		dt, err := NewDayTime(hour, minute)
		if err != nil {
			t.Fatalf("new day time: %v", err)
		}
		return dt
	}

	cases := []struct {
		name       string
		fromH, toH int
		wantHours  float64
	}{
		{name: "上午两节", fromH: 9, toH: 11, wantHours: 2},
		{name: "下午两节", fromH: 14, toH: 16, wantHours: 2},
		{name: "跨中午", fromH: 11, toH: 14, wantHours: 3},
		{name: "全天", fromH: 0, toH: 23, wantHours: 23},
		{name: "贴着午夜结束", fromH: 21, toH: 23, wantHours: 2},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			span, err := NewDayTimeRange(at(t, tc.fromH, 0), at(t, tc.toH, 0))
			if err != nil {
				t.Fatalf("NewDayTimeRange: %v", err)
			}
			if got := span.DurationHours(); got != tc.wantHours {
				t.Errorf("DurationHours = %v, want %v", got, tc.wantHours)
			}
		})
	}
}

// TestDayTimeRangeOverlaps 重叠判定：相接不算重叠，包含算重叠。
func TestDayTimeRangeOverlaps(t *testing.T) {
	span := func(t *testing.T, fromH, toH int) DayTimeRange {
		t.Helper()
		from, err := NewDayTime(fromH, 0)
		if err != nil {
			t.Fatalf("new day time: %v", err)
		}
		to, err := NewDayTime(toH, 0)
		if err != nil {
			t.Fatalf("new day time: %v", err)
		}
		got, err := NewDayTimeRange(from, to)
		if err != nil {
			t.Fatalf("new day time range: %v", err)
		}
		return got
	}

	cases := []struct {
		name  string
		aFrom int
		aTo   int
		bFrom int
		bTo   int
		want  bool
	}{
		{name: "完全重叠", aFrom: 9, aTo: 11, bFrom: 9, bTo: 11, want: true},
		{name: "部分重叠", aFrom: 9, aTo: 11, bFrom: 10, bTo: 12, want: true},
		{name: "包含", aFrom: 9, aTo: 12, bFrom: 10, bTo: 11, want: true},
		{name: "首尾相接", aFrom: 9, aTo: 11, bFrom: 11, bTo: 13, want: false},
		{name: "完全错开", aFrom: 9, aTo: 11, bFrom: 14, bTo: 16, want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a, b := span(t, tc.aFrom, tc.aTo), span(t, tc.bFrom, tc.bTo)
			if got := a.Overlaps(b); got != tc.want {
				t.Errorf("a.Overlaps(b) = %v, want %v", got, tc.want)
			}
			if got := b.Overlaps(a); got != tc.want {
				t.Errorf("b.Overlaps(a) = %v, want %v（应当对称）", got, tc.want)
			}
		})
	}
}

// TestInstantiateForDateStaysOnSameDay 模板落到具体日期上，起止仍落在同一天。
//
// 这是「槽位不允许跨天」的直接体现：跨午夜区间根本构造不出来，
// 实例化出来的 start/end 只可能差在时分上。
func TestInstantiateForDateStaysOnSameDay(t *testing.T) {
	from, err := NewDayTime(21, 0)
	if err != nil {
		t.Fatalf("new day time: %v", err)
	}
	to, err := NewDayTime(23, 59)
	if err != nil {
		t.Fatalf("new day time: %v", err)
	}
	span, err := NewDayTimeRange(from, to)
	if err != nil {
		t.Fatalf("new day time range: %v", err)
	}

	// 2026-09-15 是周二
	date := time.Date(2026, time.September, 15, 0, 0, 0, 0, time.Local)
	slot, err := NewCourseSlot("C001", time.Tuesday, span, PendingTeacherID, PendingClassroomID)
	if err != nil {
		t.Fatalf("new course slot: %v", err)
	}

	start, end, err := slot.InstantiateForDate(date, time.Local)
	if err != nil {
		t.Fatalf("InstantiateForDate: %v", err)
	}

	sy, sm, sd := start.Date()
	ey, em, ed := end.Date()
	if sy != ey || sm != em || sd != ed {
		t.Errorf("实例化后跨天了：start = %v, end = %v", start, end)
	}
	if !start.Equal(time.Date(2026, time.September, 15, 21, 0, 0, 0, time.Local)) {
		t.Errorf("start = %v, want 2026-09-15 21:00", start)
	}
	if !end.Equal(time.Date(2026, time.September, 15, 23, 59, 0, 0, time.Local)) {
		t.Errorf("end = %v, want 2026-09-15 23:59", end)
	}
}
