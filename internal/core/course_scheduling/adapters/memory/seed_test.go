package memory

import "testing"

func TestSeedDemoPopulatesEveryTable(t *testing.T) {
	d, cleanup, err := NewData(nil)
	if err != nil {
		t.Fatalf("NewData() error = %v", err)
	}
	defer cleanup()

	if err := d.SeedDemo(); err != nil {
		t.Fatalf("SeedDemo() error = %v", err)
	}

	tests := []struct {
		name string
		got  int
		want int
	}{
		{"courses", len(d.Courses()), 4},
		{"classrooms", len(d.Classrooms()), 2},
		{"students", len(d.Students()), 3},
		{"teachers", len(d.Teachers()), 3},
		{"course slots", len(d.CourseSlots()), 6},
		{"course slot changes", len(d.CourseSlotChanges()), 2},
		{"absences", len(d.Absences()), 2},
		{"enrollments", len(d.Enrollments()), 4},
		{"makeups", len(d.Makeups()), 2},
		{"qualifications", len(d.Qualifications()), 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("count = %d, want %d", tt.got, tt.want)
			}
		})
	}
}

// 反复调用 SeedDemo 应当按 ID 覆盖而不是追加。
func TestSeedDemoIsIdempotent(t *testing.T) {
	d, cleanup, err := NewData(nil)
	if err != nil {
		t.Fatalf("NewData() error = %v", err)
	}
	defer cleanup()

	for i := 1; i <= 3; i++ {
		if err := d.SeedDemo(); err != nil {
			t.Fatalf("SeedDemo() #%d error = %v", i, err)
		}
	}

	if got := len(d.Courses()); got != 4 {
		t.Errorf("courses after 3 seeds = %d, want 4", got)
	}
	if got := len(d.CourseSlots()); got != 6 {
		t.Errorf("course slots after 3 seeds = %d, want 6", got)
	}
}

// 快照按 ID 升序，分页顺序才稳定。
func TestSnapshotsAreSortedByID(t *testing.T) {
	d, _, err := NewData(nil)
	if err != nil {
		t.Fatalf("NewData() error = %v", err)
	}
	if err := d.SeedDemo(); err != nil {
		t.Fatalf("SeedDemo() error = %v", err)
	}

	courses := d.Courses()
	for i := 1; i < len(courses); i++ {
		if courses[i-1].ID() >= courses[i].ID() {
			t.Fatalf("courses not sorted: %q before %q", courses[i-1].ID(), courses[i].ID())
		}
	}

	slots := d.CourseSlots()
	for i := 1; i < len(slots); i++ {
		if slots[i-1].ID() >= slots[i].ID() {
			t.Fatalf("course slots not sorted: %d before %d", slots[i-1].ID(), slots[i].ID())
		}
	}
}

// 快照必须是新切片：改动返回值不应影响存储。
func TestSnapshotIsACopy(t *testing.T) {
	d, _, err := NewData(nil)
	if err != nil {
		t.Fatalf("NewData() error = %v", err)
	}
	if err := d.SeedDemo(); err != nil {
		t.Fatalf("SeedDemo() error = %v", err)
	}

	first := d.Courses()
	first[0] = nil

	if second := d.Courses(); second[0] == nil {
		t.Fatal("mutating the returned slice leaked into the store")
	}
}
