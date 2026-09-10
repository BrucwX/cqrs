package query

import (
	"context"
	"testing"
)

func TestCourseQueryPage(t *testing.T) {
	q := NewCourseQuery(newSeededData(t))
	ctx := context.Background()

	tests := []struct {
		name     string
		page     int
		pageSize int
		want     []string
	}{
		{"page 1", 1, 2, []string{"C001", "C002"}},
		{"page 2", 2, 2, []string{"C003", "C004"}},
		{"page past the end", 3, 2, []string{}},
		{"pageSize larger than total", 1, 100, []string{"C001", "C002", "C003", "C004"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := q.Page(ctx, tt.page, tt.pageSize)
			if err != nil {
				t.Fatalf("Page() error = %v", err)
			}
			assertIDs(t, courseIDs(got), tt.want)
		})
	}
}

func TestCourseQueryAvailableForStudent(t *testing.T) {
	q := NewCourseQuery(newSeededData(t))
	ctx := context.Background()

	// 数据背景：C001 占用 周一/周三 09:00-11:00，C002 占用 周一 14:00-16:00，
	// C003 占用 周三/周五 09:00-11:00，C004 占用 周二 09:00-11:00。
	tests := []struct {
		name      string
		studentID int64
		want      []string
	}{
		{"在学 C001 -> 自身的周一/周三 09-11 挡掉 C003", 101, []string{"C002", "C004"}},
		{"在学 C002 -> 只挡掉 C002", 102, []string{"C001", "C003", "C004"}},
		{"在学 C001（已结业的 C002 不算占用）", 103, []string{"C002", "C004"}},
		{"不存在的学员没有任何占用 -> 全部可用", 999, []string{"C001", "C002", "C003", "C004"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := q.AvailableForStudent(ctx, tt.studentID)
			if err != nil {
				t.Fatalf("AvailableForStudent() error = %v", err)
			}
			assertIDs(t, courseIDs(got), tt.want)
		})
	}
}

func TestCourseQueryAvailableForTeacher(t *testing.T) {
	q := NewCourseQuery(newSeededData(t))
	ctx := context.Background()

	tests := []struct {
		name      string
		teacherID int64
		want      []string
	}{
		{"张伟已排 C001（周一/周三 09-11）", 1, []string{"C002", "C004"}},
		{"李娜只排了 C002（周一 14-16）", 2, []string{"C001", "C003", "C004"}},
		{"王强周一/周二/周三/周五都被占", 3, []string{"C002"}},
		{"不存在的讲师没有任何占用 -> 全部可用", 999, []string{"C001", "C002", "C003", "C004"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := q.AvailableForTeacher(ctx, tt.teacherID)
			if err != nil {
				t.Fatalf("AvailableForTeacher() error = %v", err)
			}
			assertIDs(t, courseIDs(got), tt.want)
		})
	}
}

func TestCourseQueryAvailableForClassroom(t *testing.T) {
	q := NewCourseQuery(newSeededData(t))
	ctx := context.Background()

	tests := []struct {
		name        string
		classroomID string
		want        []string
	}{
		{"R101 周一/周三 09-11 已被占用", "R101", []string{"C002", "C004"}},
		{"R102 周一/周二/周三/周五 都被占满", "R102", []string{}},
		{"不存在的教室没有任何占用 -> 全部可用", "UNKNOWN", []string{"C001", "C002", "C003", "C004"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := q.AvailableForClassroom(ctx, tt.classroomID)
			if err != nil {
				t.Fatalf("AvailableForClassroom() error = %v", err)
			}
			assertIDs(t, courseIDs(got), tt.want)
		})
	}
}

// 只要求「同星期几 且 时间重叠」即算冲突：即使讲师与教室都不同也一样。
func TestCourseQueryAvailabilityIgnoresTeacherAndClassroom(t *testing.T) {
	d := newSeededData(t)
	q := NewCourseQuery(d)

	// C003 用王强 + R102，与 C001 用张伟 + R101 完全不同，
	// 但 C003 周三 09-11 与 C001 周三 09-11 撞时间 -> 对张伟而言 C003 不可用。
	got, err := q.AvailableForTeacher(context.Background(), 1)
	if err != nil {
		t.Fatalf("AvailableForTeacher() error = %v", err)
	}
	for _, c := range courseIDs(got) {
		if c == "C003" {
			t.Fatal("C003 与张伟已有的周三 09-11 撞时间，不应出现在可用列表里")
		}
	}
}
