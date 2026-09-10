package query

import (
	"slices"
	"testing"
)

func TestPaginate(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}

	tests := []struct {
		name     string
		page     int
		pageSize int
		want     []int
	}{
		{"first page", 1, 2, []int{1, 2}},
		{"middle page", 2, 2, []int{3, 4}},
		{"partial last page", 3, 2, []int{5}},
		{"page beyond the end", 9, 2, []int{}},
		{"page 0 falls back to page 1", 0, 2, []int{1, 2}},
		{"negative page falls back to page 1", -3, 2, []int{1, 2}},
		{"pageSize 0 falls back to default", 1, 0, []int{1, 2, 3, 4, 5}},
		{"negative pageSize falls back to default", 1, -1, []int{1, 2, 3, 4, 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := paginate(items, tt.page, tt.pageSize)
			if !slices.Equal(got, tt.want) {
				t.Errorf("paginate(%d, %d) = %v, want %v", tt.page, tt.pageSize, got, tt.want)
			}
			if got == nil {
				t.Error("paginate() must return a non-nil slice")
			}
		})
	}
}

func TestPaginateEmptyInput(t *testing.T) {
	got := paginate([]int{}, 1, 10)
	if got == nil {
		t.Fatal("paginate() on empty input must return a non-nil slice")
	}
	if len(got) != 0 {
		t.Errorf("len = %d, want 0", len(got))
	}
}

func TestPaginateDefaultPageSize(t *testing.T) {
	items := make([]int, defaultPageSize+5)
	for i := range items {
		items[i] = i
	}

	if got := paginate(items, 1, 0); len(got) != defaultPageSize {
		t.Errorf("page 1 size = %d, want %d", len(got), defaultPageSize)
	}
	if got := paginate(items, 2, 0); len(got) != 5 {
		t.Errorf("page 2 size = %d, want 5", len(got))
	}
}

func TestIndexByKeepsFirstDuplicate(t *testing.T) {
	type row struct {
		id   string
		note string
	}
	rows := []row{
		{"a", "first"},
		{"b", "only"},
		{"a", "second"},
	}

	got := indexBy(rows, func(r row) string { return r.id })

	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got["a"].note != "first" {
		t.Errorf("duplicate key kept %q, want %q", got["a"].note, "first")
	}
}
