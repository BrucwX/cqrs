package courseSlot

import (
	"context"
	"errors"
	"testing"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
)

// TestAssignCourse 课程在该时段没有别的排期时配置成功。
func TestAssignCourse(t *testing.T) {
	f := newFixture(t, true)
	slot := f.newSlot(t, otherCourse, noTeacher, time.Monday, 9, 11)

	if err := f.handler.AssignCourse(context.Background(), AssignCourseInput{
		SlotIDs:  []string{slot.ID()},
		CourseID: ownCourse,
	}); err != nil {
		t.Fatalf("AssignCourse: %v", err)
	}

	if got := f.slotCourseID(t, slot.ID()); got != ownCourse {
		t.Errorf("courseID = %q, want %q", got, ownCourse)
	}
}

// TestAssignCourseTimeConflict 该课程在目标时段已有别的排期 —> 拒绝且不写入。
func TestAssignCourseTimeConflict(t *testing.T) {
	f := newFixture(t, true)

	// 这门课周一 10:00-12:00 已经有一节 —— 与目标 09:00-11:00 重叠
	f.newSlot(t, ownCourse, noTeacher, time.Monday, 10, 12)
	target := f.newSlot(t, otherCourse, noTeacher, time.Monday, 9, 11)

	err := f.handler.AssignCourse(context.Background(), AssignCourseInput{
		SlotIDs:  []string{target.ID()},
		CourseID: ownCourse,
	})
	if !errors.Is(err, courseSlot.ErrCourseSlotConflict) {
		t.Errorf("err = %v, want %v", err, courseSlot.ErrCourseSlotConflict)
	}
	if got := f.slotCourseID(t, target.ID()); got != otherCourse {
		t.Errorf("被拒绝时不应写入，courseID = %q, want %q", got, otherCourse)
	}
}

// TestAssignCourseDuplicateAssignment 目标槽位本身就已经属于这门课 —> 判重复配置。
func TestAssignCourseDuplicateAssignment(t *testing.T) {
	f := newFixture(t, true)
	target := f.newSlot(t, ownCourse, noTeacher, time.Monday, 9, 11)

	err := f.handler.AssignCourse(context.Background(), AssignCourseInput{
		SlotIDs:  []string{target.ID()},
		CourseID: ownCourse,
	})
	if !errors.Is(err, courseSlot.ErrCourseSlotConflict) {
		t.Errorf("重复配置应判冲突，err = %v, want %v", err, courseSlot.ErrCourseSlotConflict)
	}
}

// TestAssignCourseNotFound 课程不存在 —> 报 not found。
func TestAssignCourseNotFound(t *testing.T) {
	f := newFixture(t, true)
	slot := f.newSlot(t, otherCourse, noTeacher, time.Monday, 9, 11)

	err := f.handler.AssignCourse(context.Background(), AssignCourseInput{
		SlotIDs:  []string{slot.ID()},
		CourseID: "no-such-course",
	})
	if !errors.Is(err, course.ErrCourseNotFound) {
		t.Errorf("err = %v, want %v", err, course.ErrCourseNotFound)
	}
}

// TestAssignCourseSlotNotFound 槽位不存在 —> 报 not found。
func TestAssignCourseSlotNotFound(t *testing.T) {
	f := newFixture(t, true)

	err := f.handler.AssignCourse(context.Background(), AssignCourseInput{
		SlotIDs:  []string{"no-such-slot"},
		CourseID: ownCourse,
	})
	if !errors.Is(err, courseSlot.ErrCourseSlotNotFound) {
		t.Errorf("err = %v, want %v", err, courseSlot.ErrCourseSlotNotFound)
	}
}
