package course

import (
	"context"
	"errors"
	"testing"
	"time"

	commandmemory "cqrs/internal/core/course_scheduling/adapters/command/memory"
	memorycmd "cqrs/internal/core/course_scheduling/adapters/command/memory/implement"
	"cqrs/internal/core/course_scheduling/adapters/memorystore"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
)

// newHandler 装配一个跑在干净内存存储上的命令处理器。
//
// 返回的是完整 store（测试要靠它塞数据和做断言），仓库拿到的则是收窄后的写侧面。
func newHandler(t *testing.T) (*Handler, *memorystore.Data) {
	t.Helper()

	store, cleanup, err := memorystore.NewData(nil)
	if err != nil {
		t.Fatalf("new data: %v", err)
	}
	t.Cleanup(cleanup)

	data := commandmemory.NewData(store)
	return NewHandler(memorycmd.NewCourseCommand(data), data), store
}

func newCourse(t *testing.T, id string) *course.Course {
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

	return course.Reconstitute(id, "ct-demo", capacity, enrollment, period)
}

// TestDeleteCourse 删除存在的课程成功，再删报 not found。
func TestDeleteCourse(t *testing.T) {
	h, d := newHandler(t)
	crs := newCourse(t, "C001")
	d.SeedCourse(crs)

	if err := h.DeleteCourse(context.Background(), crs.ID()); err != nil {
		t.Fatalf("DeleteCourse: %v", err)
	}
	if got := d.Courses(); len(got) != 0 {
		t.Errorf("课程数量 = %d, want 0", len(got))
	}

	if err := h.DeleteCourse(context.Background(), crs.ID()); !errors.Is(err, course.ErrCourseNotFound) {
		t.Errorf("重复删除 err = %v, want %v", err, course.ErrCourseNotFound)
	}
}
