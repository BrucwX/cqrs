package query

import (
	"context"
	"testing"

	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
)

func TestCourseEnrollmentQueryPage(t *testing.T) {
	q := NewCourseEnrollmentQuery(newSeededData(t))

	got, err := q.Page(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("Page() error = %v", err)
	}
	assertIDs(t, enrollmentIDs(got), []int64{1, 2, 3, 4})
}

func TestCourseEnrollmentQueryCoursesByStudentID(t *testing.T) {
	q := NewCourseEnrollmentQuery(newSeededData(t))
	ctx := context.Background()

	tests := []struct {
		name      string
		studentID int64
		want      []string
	}{
		{"陈晨只注册了 C001", 101, []string{"C001"}},
		{"刘洋注册了 C002", 102, []string{"C002"}},
		{"赵敏两条注册记录（含已结业的 C002）", 103, []string{"C001", "C002"}},
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

func TestCourseEnrollmentQueryStudentsByCourseID(t *testing.T) {
	q := NewCourseEnrollmentQuery(newSeededData(t))
	ctx := context.Background()

	tests := []struct {
		name     string
		courseID string
		want     []int64
	}{
		{"C001 有陈晨与赵敏", "C001", []int64{101, 103}},
		{"C002 有刘洋与赵敏（含已结业）", "C002", []int64{102, 103}},
		{"C003 无人注册", "C003", []int64{}},
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

// 当前实现不过滤注册状态：已结业的记录也会出现在结果里。
func TestCourseEnrollmentQueryDoesNotFilterStatus(t *testing.T) {
	d := newSeededData(t)
	q := NewCourseEnrollmentQuery(d)

	completed, err := q.CoursesByStudentID(context.Background(), 103)
	if err != nil {
		t.Fatalf("CoursesByStudentID() error = %v", err)
	}
	if len(completed) != 2 {
		t.Fatalf("赵敏 courses = %v, want 2 entries including the completed one", courseIDs(completed))
	}

	// 反证：这条注册确实已经是「已结业」而不是「在学」。
	for _, e := range d.Enrollments() {
		if e.ID() == 4 {
			if e.Status() != enrollment.StatusCompleted {
				t.Fatalf("enrollment 4 status = %v, want %v", e.Status(), enrollment.StatusCompleted)
			}
			return
		}
	}
	t.Fatal("enrollment 4 not found in seed data")
}
