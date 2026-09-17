package courseSlotChange

import (
	"errors"
	"testing"
	"time"
)

// 基准：2026-09-15（周二）14:00 起两小时。
func baseStart() time.Time {
	return time.Date(2026, time.September, 15, 14, 0, 0, 0, time.Local)
}

// TestNewTargetPlan 合法目标计划构造成功。
func TestNewTargetPlan(t *testing.T) {
	start := baseStart()
	tp, err := NewTargetPlan(start, start.Add(2*time.Hour), 2, "R102")
	if err != nil {
		t.Fatalf("NewTargetPlan: %v", err)
	}

	if !tp.TargetStartAt().Equal(start) {
		t.Errorf("start = %v, want %v", tp.TargetStartAt(), start)
	}
	if got := tp.TargetEndAt().Sub(start); got != 2*time.Hour {
		t.Errorf("duration = %v, want 2h", got)
	}
	if tp.TeacherID() != 2 {
		t.Errorf("teacherID = %d, want 2", tp.TeacherID())
	}
	if tp.ClassroomID() != "R102" {
		t.Errorf("classroomID = %q, want R102", tp.ClassroomID())
	}
}

// TestNewTargetPlanRejected 各类非法目标计划都被拒绝。
func TestNewTargetPlanRejected(t *testing.T) {
	start := baseStart()

	cases := []struct {
		name string
		from time.Time
		to   time.Time
		want error
	}{
		{
			name: "结束不晚于开始",
			from: start,
			to:   start,
			want: ErrInvalidTimeRange,
		},
		{
			name: "结束早于开始",
			from: start,
			to:   start.Add(-time.Hour),
			want: ErrInvalidTimeRange,
		},
		{
			name: "跨天（当天深夜到次日凌晨）",
			from: time.Date(2026, time.September, 15, 23, 0, 0, 0, time.Local),
			to:   time.Date(2026, time.September, 16, 1, 0, 0, 0, time.Local),
			want: ErrTargetCrossDay,
		},
		{
			name: "跨月",
			from: time.Date(2026, time.September, 30, 22, 0, 0, 0, time.Local),
			to:   time.Date(2026, time.October, 1, 1, 0, 0, 0, time.Local),
			want: ErrTargetCrossDay,
		},
		{
			name: "跨年",
			from: time.Date(2026, time.December, 31, 22, 0, 0, 0, time.Local),
			to:   time.Date(2027, time.January, 1, 1, 0, 0, 0, time.Local),
			want: ErrTargetCrossDay,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := NewTargetPlan(tc.from, tc.to, 2, "R102"); !errors.Is(err, tc.want) {
				t.Errorf("err = %v, want %v", err, tc.want)
			}
		})
	}
}

// TestNewTargetPlanSameDayBoundary 当天内的边界（00:00 起、23:59 结束）允许。
func TestNewTargetPlanSameDayBoundary(t *testing.T) {
	cases := []struct {
		name string
		from time.Time
		to   time.Time
	}{
		{
			name: "从 00:00 开始",
			from: time.Date(2026, time.September, 15, 0, 0, 0, 0, time.Local),
			to:   time.Date(2026, time.September, 15, 9, 0, 0, 0, time.Local),
		},
		{
			name: "到 23:59 结束",
			from: time.Date(2026, time.September, 15, 21, 0, 0, 0, time.Local),
			to:   time.Date(2026, time.September, 15, 23, 59, 0, 0, time.Local),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := NewTargetPlan(tc.from, tc.to, 2, "R102"); err != nil {
				t.Errorf("NewTargetPlan: %v", err)
			}
		})
	}
}

// TestNewTargetPlanRejectsBadTeacherOrClassroom 讲师与教室必填。
func TestNewTargetPlanRejectsBadTeacherOrClassroom(t *testing.T) {
	start := baseStart()

	if _, err := NewTargetPlan(start, start.Add(time.Hour), 0, "R102"); err == nil {
		t.Error("teacherID = 0 应该被拒绝")
	}
	if _, err := NewTargetPlan(start, start.Add(time.Hour), 2, ""); err == nil {
		t.Error("classroomID 为空应该被拒绝")
	}
}
