package implement

import (
	"context"
	"testing"

	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
)

func TestStudentMakeupQueryPage(t *testing.T) {
	q := NewStudentMakeupQuery(newSeededData(t))

	got, err := q.Page(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("Page() error = %v", err)
	}
	assertIDs(t, makeupIDs(got), []int64{1, 2})
}

func TestStudentMakeupQueryCoursesByStudentID(t *testing.T) {
	q := NewStudentMakeupQuery(newSeededData(t))
	ctx := context.Background()

	tests := []struct {
		name      string
		studentID int64
		want      []string
	}{
		{"陈晨在 C001 补过课", 101, []string{"C001"}},
		{"刘洋在 C002 有待审批的补课", 102, []string{"C002"}},
		{"赵敏没有补课记录", 103, []string{}},
		{"不存在的学员", 999, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := q.CoursesByStudentID(ctx, tt.studentID)
			if err != nil {
				t.Fatalf("CoursesByStudentID() error = %v", err)
			}
			assertIDs(t, courseIDs(got), tt.want)
		})
	}
}

func TestStudentMakeupQueryStudentsByCourseID(t *testing.T) {
	q := NewStudentMakeupQuery(newSeededData(t))
	ctx := context.Background()

	tests := []struct {
		name     string
		courseID string
		want     []int64
	}{
		{"C001 只有陈晨", "C001", []int64{101}},
		{"C002 只有刘洋", "C002", []int64{102}},
		{"C003 没有补课记录", "C003", []int64{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := q.StudentsByCourseID(ctx, tt.courseID)
			if err != nil {
				t.Fatalf("StudentsByCourseID() error = %v", err)
			}
			assertIDs(t, studentIDs(got), tt.want)
		})
	}
}

// 当前实现不按补课状态过滤：已预约但还没去补的也算「有补课记录」。
func TestStudentMakeupQueryDoesNotFilterStatus(t *testing.T) {
	d := newSeededData(t)

	for _, m := range d.Makeups() {
		if m.ID() != 2 {
			continue
		}
		if m.Status() != makeup.StatusBooked {
			t.Fatalf("makeup 2 status = %v, want %v", m.Status(), makeup.StatusBooked)
		}

		q := NewStudentMakeupQuery(d)
		got, err := q.CoursesByStudentID(context.Background(), 102)
		if err != nil {
			t.Fatalf("CoursesByStudentID() error = %v", err)
		}
		assertIDs(t, courseIDs(got), []string{"C002"})
		return
	}
	t.Fatal("makeup 2 not found in seed data")
}
