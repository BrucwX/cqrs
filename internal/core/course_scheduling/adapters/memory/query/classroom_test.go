package query

import (
	"context"
	"testing"

	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
)

func TestClassroomQueryPage(t *testing.T) {
	q := NewClassroomQuery(newSeededData(t))
	ctx := context.Background()

	ids := func(items []*classroom.Classroom) []string {
		out := make([]string, 0, len(items))
		for _, item := range items {
			out = append(out, item.ID())
		}
		return out
	}

	tests := []struct {
		name     string
		page     int
		pageSize int
		want     []string
	}{
		{"all", 1, 10, []string{"R101", "R102"}},
		{"page 1 of 2", 1, 1, []string{"R101"}},
		{"page 2 of 2", 2, 1, []string{"R102"}},
		{"past the end", 3, 1, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := q.Page(ctx, tt.page, tt.pageSize)
			if err != nil {
				t.Fatalf("Page() error = %v", err)
			}
			assertIDs(t, ids(got), tt.want)
		})
	}
}
