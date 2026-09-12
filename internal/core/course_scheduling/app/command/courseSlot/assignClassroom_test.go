package courseSlot

import (
	"context"
	"errors"
	"testing"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
)

// TestAssignClassroom 教室装得下且该时段空闲时安排成功。
func TestAssignClassroom(t *testing.T) {
	f := newFixture(t, true)
	slot := f.newSlot(t, ownCourse, noTeacher, time.Monday, 9, 11)

	if err := f.handler.AssignClassroom(context.Background(), AssignClassroomInput{
		SlotIDs:     []string{slot.ID()},
		ClassroomID: room,
	}); err != nil {
		t.Fatalf("AssignClassroom: %v", err)
	}

	if got := f.slotClassroomID(t, slot.ID()); got != room {
		t.Errorf("classroomID = %q, want %q", got, room)
	}
}

// TestAssignClassroomCapacityExceeded 教室装不下课程人数 —> 拒绝且不写入。
func TestAssignClassroomCapacityExceeded(t *testing.T) {
	f := newFixture(t, true)
	target := f.newSlot(t, ownCourse, noTeacher, time.Monday, 9, 11)

	err := f.handler.AssignClassroom(context.Background(), AssignClassroomInput{
		SlotIDs:     []string{target.ID()},
		ClassroomID: smallRoom, // 5 座 < 课程 20 人
	})
	if !errors.Is(err, courseSlot.ErrCourseSlotConflict) {
		t.Errorf("err = %v, want %v", err, courseSlot.ErrCourseSlotConflict)
	}
	if got := f.slotClassroomID(t, target.ID()); got != noClassroom {
		t.Errorf("被拒绝时不应写入，classroomID = %q, want %q", got, noClassroom)
	}
}

// TestAssignClassroomTimeConflict 教室在目标时段已被别的排期占用 —> 拒绝。
func TestAssignClassroomTimeConflict(t *testing.T) {
	f := newFixture(t, true)

	// 这间教室周一 10:00-12:00 已经有课 —— 与目标 09:00-11:00 重叠
	f.newSlotIn(t, otherCourse, noTeacher, room, time.Monday, 10, 12)
	target := f.newSlot(t, ownCourse, noTeacher, time.Monday, 9, 11)

	err := f.handler.AssignClassroom(context.Background(), AssignClassroomInput{
		SlotIDs:     []string{target.ID()},
		ClassroomID: room,
	})
	if !errors.Is(err, courseSlot.ErrCourseSlotConflict) {
		t.Errorf("err = %v, want %v", err, courseSlot.ErrCourseSlotConflict)
	}
}

// TestAssignClassroomDuplicateAssignment 目标槽位本身就已经排了这间教室 —> 判重复安排。
func TestAssignClassroomDuplicateAssignment(t *testing.T) {
	f := newFixture(t, true)
	target := f.newSlotIn(t, ownCourse, noTeacher, room, time.Monday, 9, 11)

	err := f.handler.AssignClassroom(context.Background(), AssignClassroomInput{
		SlotIDs:     []string{target.ID()},
		ClassroomID: room,
	})
	if !errors.Is(err, courseSlot.ErrCourseSlotConflict) {
		t.Errorf("重复安排应判冲突，err = %v, want %v", err, courseSlot.ErrCourseSlotConflict)
	}
}

// TestAssignClassroomNotFound 教室不存在 —> 报 not found。
func TestAssignClassroomNotFound(t *testing.T) {
	f := newFixture(t, true)
	slot := f.newSlot(t, ownCourse, noTeacher, time.Monday, 9, 11)

	err := f.handler.AssignClassroom(context.Background(), AssignClassroomInput{
		SlotIDs:     []string{slot.ID()},
		ClassroomID: "no-such-room",
	})
	if !errors.Is(err, classroom.ErrClassroomNotFound) {
		t.Errorf("err = %v, want %v", err, classroom.ErrClassroomNotFound)
	}
}

// TestAssignClassroomCourseNotFound 目标槽位指向的课程不存在 —> 报 not found。
func TestAssignClassroomCourseNotFound(t *testing.T) {
	f := newFixture(t, true)
	slot := f.newSlot(t, otherCourse, noTeacher, time.Monday, 9, 11)

	err := f.handler.AssignClassroom(context.Background(), AssignClassroomInput{
		SlotIDs:     []string{slot.ID()},
		ClassroomID: room,
	})
	if !errors.Is(err, course.ErrCourseNotFound) {
		t.Errorf("err = %v, want %v", err, course.ErrCourseNotFound)
	}
}
