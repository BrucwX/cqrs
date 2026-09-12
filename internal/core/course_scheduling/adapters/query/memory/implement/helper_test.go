package implement

import (
	"slices"
	"testing"

	"cqrs/internal/core/course_scheduling/adapters/memorystore"
	"cqrs/internal/core/course_scheduling/adapters/query/memory"
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

// newSeededData 返回注入了 SeedDemo 演示数据的内存存储，已收窄成只读面。
//
// 绝大多数查询测试用这个就够：数据固定，期望值稳定可复现。
func newSeededData(t *testing.T) *memory.Data {
	return memory.NewData(newSeededStore(t))
}

// newSeededStore 同上，但返回完整 store 而不收窄。
//
// 少数用例要在演示数据之上再补几条（比如给某个讲师补一条资质），
// 那就得拿到完整 store —— 收窄之后 Seed* 是调不到的。
func newSeededStore(t *testing.T) *memorystore.Data {
	t.Helper()

	store, _, err := memorystore.NewData(nil)
	if err != nil {
		t.Fatalf("memorystore.NewData() error = %v", err)
	}
	if err := store.SeedDemo(); err != nil {
		t.Fatalf("SeedDemo() error = %v", err)
	}
	return store
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
