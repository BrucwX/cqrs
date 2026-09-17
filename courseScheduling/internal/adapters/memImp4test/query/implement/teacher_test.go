package implement

import (
	"context"
	"testing"
)

func TestTeacherQueryPage(t *testing.T) {
	q := NewTeacherQuery(newSeededData(t))
	ctx := context.Background()

	tests := []struct {
		name     string
		page     int
		pageSize int
		want     []int64
	}{
		{"all", 1, 10, []int64{1, 2, 3}},
		{"page 1 of 3", 1, 1, []int64{1}},
		{"page 3 of 3", 3, 1, []int64{3}},
		{"past the end", 4, 1, []int64{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := q.PageTeachers(ctx, tt.page, tt.pageSize)
			if err != nil {
				t.Fatalf("Page() error = %v", err)
			}
			assertIDs(t, teacherIDs(got), tt.want)
		})
	}
}
