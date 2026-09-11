package courseSlot

import (
	"context"
	"errors"
	"testing"
	"time"

	"cqrs/internal/core/course_scheduling/adapters/memory"
	memorycmd "cqrs/internal/core/course_scheduling/adapters/memory/command"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

const (
	ownCourse   = "C001"
	otherCourse = "C999"
)

// fixture 是一份搭好的内存数据 + 处理器。
type fixture struct {
	handler *Handler
	data    *memory.Data
	teacher teacher.Teacher
	course  *course.Course
}

// newFixture 造一份数据：1 个课程类型、1 门课、1 位持证讲师。
func newFixture(t *testing.T, withQualification bool) *fixture {
	t.Helper()

	d, cleanup, err := memory.NewData(nil)
	if err != nil {
		t.Fatalf("new data: %v", err)
	}
	t.Cleanup(cleanup)

	contact, err := teacher.NewContactInfo("13800000000", "")
	if err != nil {
		t.Fatalf("new contact: %v", err)
	}
	teach, err := teacher.NewTeacher("李娜", "讲师", contact)
	if err != nil {
		t.Fatalf("new teacher: %v", err)
	}

	ct, err := courseType.NewCourseType("少儿编程", "")
	if err != nil {
		t.Fatalf("new course type: %v", err)
	}

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
	crs := course.NewCourse(ownCourse, ct.ID(), capacity, enrollment, period)

	d.SeedTeacher(teach)
	d.SeedCourseType(ct)
	d.SeedCourse(crs)

	if withQualification {
		q, err := qualification.NewQualification(
			teach.ID(), ct.ID(), now.AddDate(0, 0, -30), now.AddDate(0, 12, 0),
		)
		if err != nil {
			t.Fatalf("new qualification: %v", err)
		}
		d.SeedQualification(q)
	}

	return &fixture{
		handler: NewHandler(
			memorycmd.NewCourseSlotCommand(d),
			memorycmd.NewAssignTeacherRepo(d),
		),
		data:    d,
		teacher: *teach,
		course:  crs,
	}
}

// newSlot 造一条待排的槽位并塞进存储。
func (f *fixture) newSlot(t *testing.T, courseID string, teacherID int64, weekday time.Weekday, fromHour, toHour int) *courseSlot.CourseSlot {
	t.Helper()

	from, err := courseSlot.NewDayTime(fromHour, 0)
	if err != nil {
		t.Fatalf("new day time: %v", err)
	}
	to, err := courseSlot.NewDayTime(toHour, 0)
	if err != nil {
		t.Fatalf("new day time: %v", err)
	}
	span, err := courseSlot.NewDayTimeRange(from, to)
	if err != nil {
		t.Fatalf("new day time range: %v", err)
	}

	cs, err := courseSlot.NewCourseSlot(courseID, weekday, span, teacherID, "-")
	if err != nil {
		t.Fatalf("new course slot: %v", err)
	}
	f.data.SeedCourseSlot(cs)

	return cs
}

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
