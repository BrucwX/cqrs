package query

import (
	"context"
	"testing"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
)

func TestCourseSlotChangeQueryPage(t *testing.T) {
	q := NewCourseSlotChangeQuery(newSeededData(t))

	got, err := q.Page(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("Page() error = %v", err)
	}
	assertIDs(t, changeIDs(got), []int64{1, 2})
}

func TestCourseSlotChangeQueryListByCourseID(t *testing.T) {
	q := NewCourseSlotChangeQuery(newSeededData(t))

	tests := []struct {
		name     string
		courseID string
		want     []int64
	}{
		{"C001 有调课与代课两条", "C001", []int64{1, 2}},
		{"C002 没有变更单", "C002", []int64{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := q.ListByCourseID(tt.courseID)
			if err != nil {
				t.Fatalf("ListByCourseID() error = %v", err)
			}
			assertIDs(t, changeIDs(got), tt.want)
		})
	}
}

func TestCourseSlotChangeQueryReturnsSeededFields(t *testing.T) {
	q := NewCourseSlotChangeQuery(newSeededData(t))

	got, err := q.ListByCourseID("C001")
	if err != nil {
		t.Fatalf("ListByCourseID() error = %v", err)
	}

	// 第 1 条：调课，目标改到 2026-09-15 14:00 由张伟在 R102 上。
	reschedule := got[0]
	if reschedule.ChangeType() != courseSlotChange.TypeReschedule {
		t.Errorf("changeType = %v, want %v", reschedule.ChangeType(), courseSlotChange.TypeReschedule)
	}

	expectedStart := time.Date(2026, time.September, 15, 14, 0, 0, 0, time.Local)
	if start := reschedule.TargetPlan().TargetStartAt(); !start.Equal(expectedStart) {
		t.Errorf("target start = %v, want %v", start, expectedStart)
	}
	if got := reschedule.TargetPlan().ClassroomID(); got != "R102" {
		t.Errorf("target classroom = %q, want %q", got, "R102")
	}

	// 第 2 条：代课，原计划快照应保留原讲师与教室。
	substitute := got[1]
	if substitute.ChangeType() != courseSlotChange.TypeSubstitute {
		t.Errorf("changeType = %v, want %v", substitute.ChangeType(), courseSlotChange.TypeSubstitute)
	}
	if got := substitute.OriginalPlan().TeacherID(); got != 1 {
		t.Errorf("original teacher = %d, want 1", got)
	}
	if got := substitute.TargetPlan().TeacherID(); got != 2 {
		t.Errorf("target teacher = %d, want 2", got)
	}
}
