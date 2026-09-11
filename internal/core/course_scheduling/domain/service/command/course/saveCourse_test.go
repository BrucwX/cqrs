package course

import (
	"context"
	"errors"
	"testing"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

func strPtr(v string) *string { return &v }

// newCmd 造一条合法的新建课程命令。
func newCmd(t *testing.T) CourseInput {
	t.Helper()

	capacity, err := course.NewCapacity(20, 0)
	if err != nil {
		t.Fatalf("new capacity: %v", err)
	}
	now := time.Now()
	period, err := course.NewCoursePeriod(now.AddDate(0, 0, -7), now.AddDate(0, 6, 0), 16, 0)
	if err != nil {
		t.Fatalf("new period: %v", err)
	}
	enrollment := course.NewEnrollmentWindow(
		now.AddDate(0, 0, -1), now.AddDate(0, 0, 7), now.AddDate(0, 0, 14),
	)

	return CourseInput{
		CourseTypeID: strPtr("ct-0001"),
		Capacity:     &capacity,
		Enrollment:   &enrollment,
		Period:       &period,
	}
}

// TestSaveCourseCreates ID 为 nil 时新建，ID 由聚合生成。
func TestSaveCourseCreates(t *testing.T) {
	h, d := newHandler(t)

	if err := h.SaveCourse(context.Background(), newCmd(t)); err != nil {
		t.Fatalf("SaveCourse: %v", err)
	}

	courses := d.Courses()
	if len(courses) != 1 {
		t.Fatalf("课程数量 = %d, want 1", len(courses))
	}
	got := courses[0]
	if got.ID() == "" {
		t.Error("新建课程应当自动分配 ID")
	}
	if got.CourseTypeID() != "ct-0001" {
		t.Errorf("courseTypeID = %q, want ct-0001", got.CourseTypeID())
	}
	if got.Capacity().Max() != 20 {
		t.Errorf("capacity = %d, want 20", got.Capacity().Max())
	}
}

// TestSaveCourseUpdatesPartially 更新时只改命令里给出的字段。
func TestSaveCourseUpdatesPartially(t *testing.T) {
	h, d := newHandler(t)

	if err := h.SaveCourse(context.Background(), newCmd(t)); err != nil {
		t.Fatalf("SaveCourse: %v", err)
	}
	id := d.Courses()[0].ID()

	// 只给容量
	capacity, err := course.NewCapacity(50, 0)
	if err != nil {
		t.Fatalf("new capacity: %v", err)
	}
	if err := h.SaveCourse(context.Background(), CourseInput{
		ID:       &id,
		Capacity: &capacity,
	}); err != nil {
		t.Fatalf("SaveCourse: %v", err)
	}

	got, _ := d.CourseByID(id)
	if got.Capacity().Max() != 50 {
		t.Errorf("capacity = %d, want 50", got.Capacity().Max())
	}
	if got.CourseTypeID() != "ct-0001" {
		t.Errorf("courseTypeID = %q, want ct-0001（未给出的字段应保持原值）", got.CourseTypeID())
	}
}

// TestSaveCourseUpdateMissing 更新不存在的课程报 not found，不会顺手新建。
func TestSaveCourseUpdateMissing(t *testing.T) {
	h, d := newHandler(t)

	err := h.SaveCourse(context.Background(), CourseInput{ID: strPtr("no-such-course")})
	if !errors.Is(err, repo.ErrCourseNotFound) {
		t.Errorf("err = %v, want %v", err, repo.ErrCourseNotFound)
	}
	if courses := d.Courses(); len(courses) != 0 {
		t.Errorf("被拒绝时不应新建，课程数量 = %d", len(courses))
	}
}

// TestSaveCourseInputRejected 新建必填项缺失 / 非法课程类型都被拒绝。
func TestSaveCourseInputRejected(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(cmd *CourseInput)
		want   error
	}{
		{name: "缺课程类型", mutate: func(cmd *CourseInput) { cmd.CourseTypeID = nil }, want: course.ErrCourseTypeRequired},
		{name: "课程类型为空串", mutate: func(cmd *CourseInput) { cmd.CourseTypeID = strPtr("") }, want: course.ErrCourseTypeRequired},
		{name: "缺容量", mutate: func(cmd *CourseInput) { cmd.Capacity = nil }, want: course.ErrCapacityRequired},
		{name: "缺选课窗口", mutate: func(cmd *CourseInput) { cmd.Enrollment = nil }, want: course.ErrEnrollmentRequired},
		{name: "缺课程周期", mutate: func(cmd *CourseInput) { cmd.Period = nil }, want: course.ErrPeriodRequired},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h, d := newHandler(t)
			cmd := newCmd(t)
			tc.mutate(&cmd)

			err := h.SaveCourse(context.Background(), cmd)
			if !errors.Is(err, tc.want) {
				t.Errorf("err = %v, want %v", err, tc.want)
			}
			if courses := d.Courses(); len(courses) != 0 {
				t.Errorf("被拒绝时不应写入，课程数量 = %d", len(courses))
			}
		})
	}
}
