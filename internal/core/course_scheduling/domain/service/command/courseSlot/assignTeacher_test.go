package courseSlot

import (
	"context"
	"errors"
	"testing"
	"time"

	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// TestAssignTeacher 有资质且时间不冲突时排课成功。
func TestAssignTeacher(t *testing.T) {
	f := newFixture(t, true)
	slot := f.newSlot(t, ownCourse, -1, time.Monday, 9, 11)

	if err := f.handler.AssignTeacher(context.Background(), AssignTeacherInput{
		SlotIDs:   []string{slot.ID()},
		TeacherID: f.teacher.ID(),
	}); err != nil {
		t.Fatalf("AssignTeacher: %v", err)
	}

	saved, ok := f.data.CourseSlotByID(slot.ID())
	if !ok {
		t.Fatal("slot not found after assign")
	}
	if saved.TeacherID() != f.teacher.ID() {
		t.Errorf("teacherID = %d, want %d", saved.TeacherID(), f.teacher.ID())
	}
}

// TestAssignTeacherWithoutQualification 讲师没有该课程类型的资质 —> 拒绝且不写入。
func TestAssignTeacherWithoutQualification(t *testing.T) {
	f := newFixture(t, false)
	slot := f.newSlot(t, ownCourse, -1, time.Monday, 9, 11)

	err := f.handler.AssignTeacher(context.Background(), AssignTeacherInput{
		SlotIDs:   []string{slot.ID()},
		TeacherID: f.teacher.ID(),
	})
	if err == nil {
		t.Fatal("期望被拒绝，实际成功")
	}
	if !errors.Is(err, repo.ErrCourseSlotConflict) {
		t.Errorf("err = %v, want %v", err, repo.ErrCourseSlotConflict)
	}

	saved, _ := f.data.CourseSlotByID(slot.ID())
	if saved.TeacherID() != -1 {
		t.Errorf("被拒绝时不应写入，teacherID = %d, want -1", saved.TeacherID())
	}
}

// TestAssignTeacherTimeConflict 讲师在目标时段已有别的课 —> 拒绝。
func TestAssignTeacherTimeConflict(t *testing.T) {
	f := newFixture(t, true)

	// 讲师周一 10:00-12:00 已在别的课上 —— 与目标 09:00-11:00 重叠
	f.newSlot(t, otherCourse, f.teacher.ID(), time.Monday, 10, 12)
	target := f.newSlot(t, ownCourse, -1, time.Monday, 9, 11)

	err := f.handler.AssignTeacher(context.Background(), AssignTeacherInput{
		SlotIDs:   []string{target.ID()},
		TeacherID: f.teacher.ID(),
	})
	if !errors.Is(err, repo.ErrCourseSlotConflict) {
		t.Errorf("err = %v, want %v", err, repo.ErrCourseSlotConflict)
	}
}

// TestAssignTeacherOtherSlotOverlaps 该讲师另一条槽位与目标重叠 —> 拒绝。
func TestAssignTeacherOtherSlotOverlaps(t *testing.T) {
	f := newFixture(t, true)

	// 同一位讲师、同一时段的另一条槽位，跟目标会重叠 —— 应该判冲突
	f.newSlot(t, otherCourse, f.teacher.ID(), time.Monday, 9, 11)
	target := f.newSlot(t, ownCourse, f.teacher.ID(), time.Monday, 9, 11)

	err := f.handler.AssignTeacher(context.Background(), AssignTeacherInput{
		SlotIDs:   []string{target.ID()},
		TeacherID: f.teacher.ID(),
	})
	if !errors.Is(err, repo.ErrCourseSlotConflict) {
		t.Errorf("与别的槽位重叠时应判冲突，err = %v", err)
	}
}

// TestAssignTeacherDuplicateAssignment 目标槽位本身就已经是该讲师的 —> 判重复安排，拒绝。
func TestAssignTeacherDuplicateAssignment(t *testing.T) {
	f := newFixture(t, true)
	target := f.newSlot(t, ownCourse, f.teacher.ID(), time.Monday, 9, 11)

	err := f.handler.AssignTeacher(context.Background(), AssignTeacherInput{
		SlotIDs:   []string{target.ID()},
		TeacherID: f.teacher.ID(),
	})
	if !errors.Is(err, repo.ErrCourseSlotConflict) {
		t.Errorf("重复安排应判冲突，err = %v, want %v", err, repo.ErrCourseSlotConflict)
	}
}

// TestAssignTeacherNotFound 讲师不存在 —> 报 not found。
func TestAssignTeacherNotFound(t *testing.T) {
	f := newFixture(t, true)
	slot := f.newSlot(t, ownCourse, -1, time.Monday, 9, 11)

	err := f.handler.AssignTeacher(context.Background(), AssignTeacherInput{
		SlotIDs:   []string{slot.ID()},
		TeacherID: 999999,
	})
	if !errors.Is(err, repo.ErrTeacherNotFound) {
		t.Errorf("err = %v, want %v", err, repo.ErrTeacherNotFound)
	}
}

// TestAssignTeacherSlotNotFound 槽位不存在 —> 报 not found。
func TestAssignTeacherSlotNotFound(t *testing.T) {
	f := newFixture(t, true)

	err := f.handler.AssignTeacher(context.Background(), AssignTeacherInput{
		SlotIDs:   []string{"no-such-slot"},
		TeacherID: f.teacher.ID(),
	})
	if !errors.Is(err, repo.ErrCourseSlotNotFound) {
		t.Errorf("err = %v, want %v", err, repo.ErrCourseSlotNotFound)
	}
}
