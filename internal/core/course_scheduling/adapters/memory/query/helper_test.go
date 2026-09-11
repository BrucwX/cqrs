package query

import (
	"slices"
	"testing"

	"cqrs/internal/core/course_scheduling/adapters/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
)

// newSeededData 返回注入了 SeedDemo 演示数据的内存存储。
// 所有查询测试共用这一份固定数据，所以期望值稳定可复现。
func newSeededData(t *testing.T) *memory.Data {
	t.Helper()

	d, _, err := memory.NewData(nil)
	if err != nil {
		t.Fatalf("memory.NewData() error = %v", err)
	}
	if err := d.SeedDemo(); err != nil {
		t.Fatalf("SeedDemo() error = %v", err)
	}
	return d
}

// assertIDs 断言两个 ID 序列完全一致（顺序敏感）。
func assertIDs[T comparable](t *testing.T, got, want []T) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func courseIDs(items []*course.Course) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, item.ID())
	}
	return out
}

func courseTypeIDs(items []*courseType.CourseType) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, item.ID())
	}
	return out
}

func studentIDs(items []*student.Student) []int64 {
	out := make([]int64, 0, len(items))
	for _, item := range items {
		out = append(out, item.ID())
	}
	return out
}

func teacherIDs(items []*teacher.Teacher) []int64 {
	out := make([]int64, 0, len(items))
	for _, item := range items {
		out = append(out, item.ID())
	}
	return out
}

func slotIDs(items []*courseSlot.CourseSlot) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, item.ID())
	}
	return out
}

func changeIDs(items []*courseSlotChange.CourseSlotChange) []int64 {
	out := make([]int64, 0, len(items))
	for _, item := range items {
		out = append(out, item.ID())
	}
	return out
}

func absenceIDs(items []*absence.AbsenceRecord) []int64 {
	out := make([]int64, 0, len(items))
	for _, item := range items {
		out = append(out, item.ID())
	}
	return out
}

func enrollmentIDs(items []*enrollment.CourseEnrollment) []int64 {
	out := make([]int64, 0, len(items))
	for _, item := range items {
		out = append(out, item.ID())
	}
	return out
}

func makeupIDs(items []*makeup.StudentMakeup) []int64 {
	out := make([]int64, 0, len(items))
	for _, item := range items {
		out = append(out, item.ID())
	}
	return out
}

func qualificationIDs(items []*qualification.Qualification) []int64 {
	out := make([]int64, 0, len(items))
	for _, item := range items {
		out = append(out, item.ID())
	}
	return out
}
