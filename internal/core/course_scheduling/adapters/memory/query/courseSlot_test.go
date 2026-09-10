package query

import (
	"context"
	"testing"
	"time"
)

func TestCourseSlotQueryPage(t *testing.T) {
	q := NewCourseSlotQuery(newSeededData(t))

	got, err := q.Page(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("Page() error = %v", err)
	}
	assertIDs(t, slotIDs(got), []int64{1, 2, 3, 4, 5, 6})
}

func TestCourseSlotQueryListByCourseID(t *testing.T) {
	q := NewCourseSlotQuery(newSeededData(t))

	tests := []struct {
		name     string
		courseID string
		want     []int64
	}{
		{"C001 两次课", "C001", []int64{1, 2}},
		{"C002 一次课", "C002", []int64{3}},
		{"C003 两次课", "C003", []int64{4, 5}},
		{"C004 一次课", "C004", []int64{6}},
		{"不存在的课程", "C999", []int64{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := q.ListByCourseID(tt.courseID)
			if err != nil {
				t.Fatalf("ListByCourseID() error = %v", err)
			}
			assertIDs(t, slotIDs(got), tt.want)
		})
	}
}

// 顺带校验演示数据的排期字段确实被正确读出来。
func TestCourseSlotQueryReturnsSeededFields(t *testing.T) {
	q := NewCourseSlotQuery(newSeededData(t))

	got, err := q.ListByCourseID("C001")
	if err != nil {
		t.Fatalf("ListByCourseID() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}

	first := got[0]
	if first.Weekday() != time.Monday {
		t.Errorf("weekday = %v, want %v", first.Weekday(), time.Monday)
	}
	if got := first.TeacherID(); got != 1 {
		t.Errorf("teacherID = %d, want 1", got)
	}
	if got := first.ClassroomID(); got != "R101" {
		t.Errorf("classroomID = %q, want %q", got, "R101")
	}
	if got := first.TimeRange().DurationHours(); got != 2 {
		t.Errorf("duration = %v hours, want 2", got)
	}
}
