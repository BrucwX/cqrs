package implement

import (
	"context"
	"testing"
)

func TestAbsenceRecordQueryPage(t *testing.T) {
	q := NewAbsenceRecordQuery(newSeededData(t))

	got, err := q.Page(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("Page() error = %v", err)
	}
	assertIDs(t, absenceIDs(got), []int64{1, 2})
}

func TestAbsenceRecordQueryListByStudentID(t *testing.T) {
	q := NewAbsenceRecordQuery(newSeededData(t))

	tests := []struct {
		name      string
		studentID int64
		want      []int64
	}{
		{"陈晨的事假", 101, []int64{1}},
		{"刘洋的旷课", 102, []int64{2}},
		{"赵敏没有缺勤记录", 103, []int64{}},
		{"不存在的学员", 999, []int64{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := q.ListByStudentID(tt.studentID)
			if err != nil {
				t.Fatalf("ListByStudentID() error = %v", err)
			}
			assertIDs(t, absenceIDs(got), tt.want)
		})
	}
}

func TestAbsenceRecordQueryListByCourseID(t *testing.T) {
	q := NewAbsenceRecordQuery(newSeededData(t))

	tests := []struct {
		name     string
		courseID string
		want     []int64
	}{
		{"C001", "C001", []int64{1}},
		{"C002", "C002", []int64{2}},
		{"C003 没有缺勤记录", "C003", []int64{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := q.ListByCourseID(tt.courseID)
			if err != nil {
				t.Fatalf("ListByCourseID() error = %v", err)
			}
			assertIDs(t, absenceIDs(got), tt.want)
		})
	}
}
