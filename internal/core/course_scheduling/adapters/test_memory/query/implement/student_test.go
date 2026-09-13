package implement

import (
	"context"
	"testing"
)

func TestStudentQueryPage(t *testing.T) {
	q := NewStudentQuery(newSeededData(t))
	ctx := context.Background()

	tests := []struct {
		name     string
		page     int
		pageSize int
		want     []int64
	}{
		{"all", 1, 10, []int64{101, 102, 103}},
		{"page 1 of 2", 1, 2, []int64{101, 102}},
		{"page 2 of 2", 2, 2, []int64{103}},
		{"past the end", 5, 2, []int64{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := q.Page(ctx, tt.page, tt.pageSize)
			if err != nil {
				t.Fatalf("Page() error = %v", err)
			}
			assertIDs(t, studentIDs(got), tt.want)
		})
	}
}
