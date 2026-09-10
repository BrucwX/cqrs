package query

import (
	"context"
	"testing"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
)

func TestQualificationQueryPage(t *testing.T) {
	q := NewQualificationQuery(newSeededData(t))

	got, err := q.Page(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("Page() error = %v", err)
	}
	assertIDs(t, qualificationIDs(got), []int64{1, 2, 3, 4, 5})
}

func TestQualificationQueryCoursesByTeacherID(t *testing.T) {
	q := NewQualificationQuery(newSeededData(t))
	ctx := context.Background()

	tests := []struct {
		name      string
		teacherID int64
		want      []string
	}{
		{"张伟：C001 有效 + C004 已吊销", 1, []string{"C001", "C004"}},
		{"李娜：C002", 2, []string{"C002"}},
		{"王强：C003 + C004", 3, []string{"C003", "C004"}},
		{"没录过资质的讲师", 999, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := q.CoursesByTeacherID(ctx, tt.teacherID)
			if err != nil {
				t.Fatalf("CoursesByTeacherID() error = %v", err)
			}
			assertIDs(t, courseIDs(got), tt.want)
		})
	}
}

func TestQualificationQueryTeachersByCourseID(t *testing.T) {
	q := NewQualificationQuery(newSeededData(t))
	ctx := context.Background()

	tests := []struct {
		name     string
		courseID string
		want     []int64
	}{
		{"C001 只有张伟", "C001", []int64{1}},
		{"C002 只有李娜", "C002", []int64{2}},
		{"C003 只有王强", "C003", []int64{3}},
		{"C004 有王强与张伟（按资质 ID 顺序）", "C004", []int64{3, 1}},
		{"没录过资质的课程", "C999", []int64{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := q.TeachersByCourseID(ctx, tt.courseID)
			if err != nil {
				t.Fatalf("TeachersByCourseID() error = %v", err)
			}
			assertIDs(t, teacherIDs(got), tt.want)
		})
	}
}

// 当前实现不过滤资质状态：已吊销的资质仍会出现在结果里。
func TestQualificationQueryDoesNotFilterStatus(t *testing.T) {
	d := newSeededData(t)

	var revoked bool
	for _, item := range d.Qualifications() {
		if item.ID() == 5 {
			if item.Status() != qualification.StatusRevoked {
				t.Fatalf("qualification 5 status = %v, want %v", item.Status(), qualification.StatusRevoked)
			}
			revoked = true
		}
	}
	if !revoked {
		t.Fatal("qualification 5 not found in seed data")
	}

	q := NewQualificationQuery(d)
	got, err := q.CoursesByTeacherID(context.Background(), 1)
	if err != nil {
		t.Fatalf("CoursesByTeacherID() error = %v", err)
	}
	assertIDs(t, courseIDs(got), []string{"C001", "C004"})
}

// 同一个 (teacher, course) 只应出现一次，即使有多条资质记录。
func TestQualificationQueryDeduplicatesCourses(t *testing.T) {
	d := newSeededData(t)

	extra, err := qualification.NewQualification(
		6, 2, "C002", time.Now(), time.Now().AddDate(1, 0, 0),
	)
	if err != nil {
		t.Fatalf("NewQualification() error = %v", err)
	}
	d.SeedQualification(extra)

	q := NewQualificationQuery(d)
	got, err := q.CoursesByTeacherID(context.Background(), 2)
	if err != nil {
		t.Fatalf("CoursesByTeacherID() error = %v", err)
	}
	assertIDs(t, courseIDs(got), []string{"C002"})
}
