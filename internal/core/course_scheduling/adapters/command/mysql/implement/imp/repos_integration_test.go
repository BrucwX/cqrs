package imp

import (
	"context"
	"testing"
	"time"

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

// 写侧命令仓库的集成测试：对着独立测试库（见 newCommandData）做一次写入→读回
// 的往返，验证 SQL 列名与转换函数都对得上。连不上 MySQL 就跳过。

var (
	itNow   = time.Date(2026, 9, 12, 10, 0, 0, 0, time.UTC)
	itStart = time.Date(2026, 9, 16, 16, 0, 0, 0, time.UTC)
	itEnd   = time.Date(2026, 9, 16, 18, 0, 0, 0, time.UTC)
)

func TestReposRoundTrip(t *testing.T) {
	t.Run("student", func(t *testing.T) {
		c := NewStudentImp(newCommandData(t))
		contact, _ := student.NewContactInfo("13800009001", "it@example.com")
		do := student.Reconstitute(9900001, "集成学员", student.TypeExternal, contact, 1, itNow, itNow)

		if err := c.Create(context.Background(), do); err != nil {
			t.Fatalf("Create: %v", err)
		}
		t.Cleanup(func() { _ = c.Delete(context.Background(), do.ID()) })

		got, err := c.MustGet(context.Background(), do.ID())
		if err != nil {
			t.Fatalf("MustGet: %v", err)
		}
		if got.Name() != "集成学员" || got.Contact().Phone() != "13800009001" {
			t.Errorf("roundtrip mismatch: %+v", got)
		}
	})

	t.Run("teacher", func(t *testing.T) {
		c := NewTeacherImp(newCommandData(t))
		contact, _ := teacher.NewContactInfo("13900009002", "")
		do := teacher.Reconstitute(9900002, 9900001, "集成讲师", "金牌", contact, 1, itNow, itNow)

		if err := c.Create(context.Background(), do); err != nil {
			t.Fatalf("Create: %v", err)
		}
		t.Cleanup(func() { _ = c.Delete(context.Background(), do.ID()) })

		got, err := c.MustGet(context.Background(), do.ID())
		if err != nil {
			t.Fatalf("MustGet: %v", err)
		}
		if got.Name() != "集成讲师" || got.StudentID() != 9900001 {
			t.Errorf("roundtrip mismatch: %+v", got)
		}
	})

	t.Run("course", func(t *testing.T) {
		c := NewCourseImp(newCommandData(t))
		capacity, _ := course.NewCapacity(30, 5)
		period, _ := course.NewCoursePeriod(itNow, itEnd, 40, 4)
		window := course.NewEnrollmentWindow(itNow, itEnd, itStart)
		do := course.Reconstitute("it-course-1", "it-ct-1", capacity, window, period)

		if err := c.Create(context.Background(), do); err != nil {
			t.Fatalf("Create: %v", err)
		}
		t.Cleanup(func() { _ = c.Delete(context.Background(), do.ID()) })

		got, err := c.MustGet(context.Background(), do.ID())
		if err != nil {
			t.Fatalf("MustGet: %v", err)
		}
		if got.Capacity().Max() != 30 || got.CourseTypeID() != "it-ct-1" {
			t.Errorf("roundtrip mismatch: %+v", got)
		}
	})

	t.Run("course_slot", func(t *testing.T) {
		c := NewCourseSlotImp(newCommandData(t))
		s, _ := courseSlot.NewDayTime(16, 0)
		e, _ := courseSlot.NewDayTime(18, 30)
		tr, _ := courseSlot.NewDayTimeRange(s, e)
		do := courseSlot.Reconstitute("it-slot-1", "it-course-1", time.Wednesday, tr, -1, "", itNow, itNow)

		if err := c.Save(context.Background(), do); err != nil {
			t.Fatalf("Save: %v", err)
		}
		t.Cleanup(func() { _ = c.Delete(context.Background(), do.ID()) })

		got, err := c.MustGet(context.Background(), do.ID())
		if err != nil {
			t.Fatalf("MustGet: %v", err)
		}
		if got.CourseID() != "it-course-1" || got.TimeRange().Start().String() != "16:00" {
			t.Errorf("roundtrip mismatch: %+v", got)
		}

		if err := c.AssignTeacher(context.Background(), []string{do.ID()}, 9900002); err != nil {
			t.Fatalf("AssignTeacher: %v", err)
		}
		got, err = c.MustGet(context.Background(), do.ID())
		if err != nil {
			t.Fatalf("MustGet(after assign): %v", err)
		}
		if got.TeacherID() != 9900002 {
			t.Errorf("TeacherID = %d, want 9900002", got.TeacherID())
		}
	})

	t.Run("course_slot_change", func(t *testing.T) {
		c := NewCourseSlotChangeImp(newCommandData(t))
		original := courseSlotChange.NewOriginalPlan("it-slot-1", itStart, 9900002, "R101", "16:00", "18:00")
		target, _ := courseSlotChange.NewTargetPlan(itStart, itEnd, 9900002, "R102")
		do := courseSlotChange.Reconstitute(9900003, "it-course-1", 9900001, 2, original, target, "集成换课", itNow, itNow)

		if err := c.Save(context.Background(), do); err != nil {
			t.Fatalf("Save: %v", err)
		}
		t.Cleanup(func() { _ = c.Delete(context.Background(), do.ID()) })

		got, err := c.MustGet(context.Background(), do.ID())
		if err != nil {
			t.Fatalf("MustGet: %v", err)
		}
		if got.CourseID() != "it-course-1" || got.Reason() != "集成换课" {
			t.Errorf("roundtrip mismatch: %+v", got)
		}
	})

	t.Run("enrollment", func(t *testing.T) {
		c := NewCourseEnrollmentImp(newCommandData(t))
		do := enrollment.Reconstitute(9900004, 9900001, "it-course-1", 1, itNow, time.Time{}, time.Time{}, itNow)

		if err := c.Save(context.Background(), do); err != nil {
			t.Fatalf("Save: %v", err)
		}
		t.Cleanup(func() { _ = c.Delete(context.Background(), do.ID()) })

		got, err := c.MustGet(context.Background(), do.ID())
		if err != nil {
			t.Fatalf("MustGet: %v", err)
		}
		if got.StudentID() != 9900001 || !got.IsActive() {
			t.Errorf("roundtrip mismatch: %+v", got)
		}

		list, err := c.GetEnrollments(context.Background(), 9900001)
		if err != nil || len(list) != 1 {
			t.Fatalf("GetEnrollments = %v, %v; want 1 row", list, err)
		}
	})

	t.Run("absence", func(t *testing.T) {
		c := NewAbsenceRecordImp(newCommandData(t))
		do := absence.Reconstitute(9900005, 9900001, "it-course-1", "it-slot-1", itStart, 2, 3, "集成缺勤", itNow, itNow)

		if err := c.Save(context.Background(), do); err != nil {
			t.Fatalf("Save: %v", err)
		}
		t.Cleanup(func() { _ = c.Delete(context.Background(), do.ID()) })

		got, err := c.MustGet(context.Background(), do.ID())
		if err != nil {
			t.Fatalf("MustGet: %v", err)
		}
		if got.StudentID() != 9900001 || got.Reason() != "集成缺勤" {
			t.Errorf("roundtrip mismatch: %+v", got)
		}
	})

	t.Run("makeup", func(t *testing.T) {
		c := NewStudentMakeupImp(newCommandData(t))
		do := makeup.Reconstitute(9900006, 9900001, "it-course-1", "it-slot-1", itStart, "it-slot-1", itEnd, 2, 1, itEnd, itNow, itNow)

		if err := c.Save(context.Background(), do); err != nil {
			t.Fatalf("Save: %v", err)
		}
		t.Cleanup(func() { _ = c.Delete(context.Background(), do.ID()) })

		got, err := c.MustGet(context.Background(), do.ID())
		if err != nil {
			t.Fatalf("MustGet: %v", err)
		}
		if got.StudentID() != 9900001 || got.TargetSlotID() != "it-slot-1" {
			t.Errorf("roundtrip mismatch: %+v", got)
		}
	})

	t.Run("qualification", func(t *testing.T) {
		c := NewQualificationImp(newCommandData(t))
		do := qualification.Reconstitute(9900007, 9900002, "it-ct-1", itNow, 1, itNow)

		if err := c.GrantQualification(context.Background(), do); err != nil {
			t.Fatalf("GrantQualification: %v", err)
		}
		t.Cleanup(func() { _ = c.Delete(context.Background(), do.ID()) })

		got, err := c.MustGet(context.Background(), do.ID())
		if err != nil {
			t.Fatalf("MustGet: %v", err)
		}
		if got.TeacherID() != 9900002 || got.CourseTypeID() != "it-ct-1" {
			t.Errorf("roundtrip mismatch: %+v", got)
		}
	})

	t.Run("course_type", func(t *testing.T) {
		data := newCommandData(t)
		c := NewCourseTypeImp(data)
		do := courseType.Reconstitute("it-ct-1", "集成类型", "desc", itNow, itNow)

		// 课程类型接口没有写方法，这里用裸 SQL 造一行再读回。
		if _, err := data.Conn(context.Background()).ExecContext(context.Background(),
			`INSERT INTO course_type (id, name, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
			do.ID(), do.Name(), do.Description(), do.CreatedAt(), do.UpdatedAt()); err != nil {
			t.Fatalf("seed course_type: %v", err)
		}
		t.Cleanup(func() {
			_, _ = data.Conn(context.Background()).ExecContext(context.Background(),
				`DELETE FROM course_type WHERE id = ?`, do.ID())
		})

		got, err := c.MustGet(context.Background(), do.ID())
		if err != nil {
			t.Fatalf("MustGet: %v", err)
		}
		if got.Name() != "集成类型" {
			t.Errorf("roundtrip mismatch: %+v", got)
		}
	})
}
