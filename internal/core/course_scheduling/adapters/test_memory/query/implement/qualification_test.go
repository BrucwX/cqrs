package implement

import (
	"context"
	"testing"
	"time"

	"cqrs/internal/core/course_scheduling/adapters/test_memory/query"
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

func TestQualificationQueryCourseTypesByTeacherID(t *testing.T) {
	q := NewQualificationQuery(newSeededData(t))
	ctx := context.Background()

	// 演示数据：ct-0001 少儿编程、ct-0002 成人英语；C001/C002 属前者，C003/C004 属后者。
	tests := []struct {
		name      string
		teacherID int64
		want      []string
	}{
		{"张伟：少儿编程 + 成人英语", 1, []string{"ct-0001", "ct-0002"}},
		{"李娜：只有少儿编程", 2, []string{"ct-0001"}},
		{"王强：成人英语 + 少儿编程（后者已吊销，仍返回）", 3, []string{"ct-0002", "ct-0001"}},
		{"没录过资质的讲师", 999, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := q.CourseTypesByTeacherID(ctx, tt.teacherID)
			if err != nil {
				t.Fatalf("CourseTypesByTeacherID() error = %v", err)
			}
			assertIDs(t, courseTypeIDs(got), tt.want)
		})
	}
}

func TestQualificationQueryTeachersByCourseTypeID(t *testing.T) {
	q := NewQualificationQuery(newSeededData(t))
	ctx := context.Background()

	tests := []struct {
		name         string
		courseTypeID string
		want         []int64
	}{
		{"少儿编程：张伟、李娜、王强（按资质 ID 顺序）", "ct-0001", []int64{1, 2, 3}},
		{"成人英语：王强、张伟", "ct-0002", []int64{3, 1}},
		{"没录过资质的课程类型", "ct-9999", []int64{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := q.TeachersByCourseTypeID(ctx, tt.courseTypeID)
			if err != nil {
				t.Fatalf("TeachersByCourseTypeID() error = %v", err)
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
	got, err := q.CourseTypesByTeacherID(context.Background(), 3)
	if err != nil {
		t.Fatalf("CourseTypesByTeacherID() error = %v", err)
	}
	assertIDs(t, courseTypeIDs(got), []string{"ct-0002", "ct-0001"})
}

// 同一个 (teacher, courseType) 只应出现一次，即使有多条资质记录。
func TestQualificationQueryDeduplicatesCourseTypes(t *testing.T) {
	store := newSeededStore(t)

	extra := qualification.Reconstitute(
		999, 2, "ct-0001", time.Now(),
		qualification.StatusActive, time.Now(),
	)
	store.SeedQualification(extra)

	q := NewQualificationQuery(query.NewData(store))
	got, err := q.CourseTypesByTeacherID(context.Background(), 2)
	if err != nil {
		t.Fatalf("CourseTypesByTeacherID() error = %v", err)
	}
	assertIDs(t, courseTypeIDs(got), []string{"ct-0001"})
}
