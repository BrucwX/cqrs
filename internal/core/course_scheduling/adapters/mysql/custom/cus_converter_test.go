package custom

import (
	"testing"
	"time"

	"cqrs/internal/core/course_scheduling/adapters/mysql/custom/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
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

var (
	now   = time.Date(2026, 9, 12, 10, 0, 0, 0, time.UTC)
	start = time.Date(2026, 9, 16, 16, 0, 0, 0, time.UTC)
	end   = time.Date(2026, 9, 16, 18, 0, 0, 0, time.UTC)
)

// conv 就是那个唯一的转换器，所有用例共用。
var conv = Converter{}

func eq(t *testing.T, name string, got, want any) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %v want %v", name, got, want)
	}
}

// TestRoundTrip 11 个聚合各走一遍 DO -> PO -> DO，确认字段原样回来。
func TestRoundTrip(t *testing.T) {
	t.Run("courseType", func(t *testing.T) {
		do := courseType.Reconstitute("ct-1", "少儿编程", "6-12 岁", now, now)
		back, err := conv.CourseTypeToDO(conv.CourseTypeToPO(do))
		if err != nil {
			t.Fatal(err)
		}
		eq(t, "ID", back.ID(), do.ID())
		eq(t, "Description", back.Description(), do.Description())
	})

	t.Run("classroom", func(t *testing.T) {
		loc, err := classroom.NewLocation("R 座", 1, "101")
		if err != nil {
			t.Fatal(err)
		}
		do := classroom.Reconstitute("cl-1", loc, 30, 7, 2)

		po := conv.ClassroomToPO(do)
		eq(t, "Floor", po.Floor, 1)
		eq(t, "Allocated", po.Allocated, 7)

		back, err := conv.ClassroomToDO(po)
		if err != nil {
			t.Fatal(err)
		}
		eq(t, "Location", back.Location().FullName(), "R 座-1F-101")
		eq(t, "Status", back.Status(), classroom.Status(2))
	})

	t.Run("student", func(t *testing.T) {
		contact, err := student.NewContactInfo("13800000001", "a@b.com")
		if err != nil {
			t.Fatal(err)
		}
		do := student.Reconstitute(101, "陈晨", student.TypeInternal, contact, 1, now, now)

		back, err := conv.StudentToDO(conv.StudentToPO(do))
		if err != nil {
			t.Fatal(err)
		}
		eq(t, "Phone", back.Contact().Phone(), "13800000001")
		eq(t, "StudentType", back.StudentType(), student.TypeInternal)
		eq(t, "UpdatedAt", back.UpdatedAt(), now)
	})

	t.Run("teacher", func(t *testing.T) {
		contact, err := teacher.NewContactInfo("13900000002", "")
		if err != nil {
			t.Fatal(err)
		}
		do := teacher.Reconstitute(1, 201, "张伟", "金牌讲师", contact, 2, now, now)

		back, err := conv.TeacherToDO(conv.TeacherToPO(do))
		if err != nil {
			t.Fatal(err)
		}
		eq(t, "StudentID", back.StudentID(), int64(201))
		eq(t, "Title", back.Title(), "金牌讲师")
		eq(t, "Status", back.Status(), teacher.Status(2))
	})

	t.Run("course", func(t *testing.T) {
		capacity, err := course.NewCapacity(30, 5)
		if err != nil {
			t.Fatal(err)
		}
		period, err := course.NewCoursePeriod(now, end, 40, 4)
		if err != nil {
			t.Fatal(err)
		}
		window := course.NewEnrollmentWindow(now, end, start)
		do := course.Reconstitute("c-1", "ct-1", capacity, window, period)

		po := conv.CourseToPO(do)
		eq(t, "DropDeadline", po.DropDeadline, start)
		eq(t, "TotalHours", po.TotalHours, 40)

		back, err := conv.CourseToDO(po)
		if err != nil {
			t.Fatal(err)
		}
		eq(t, "Enrollment.EndAt", back.Enrollment().EndAt(), end)
		eq(t, "Period.StartAt", back.Period().StartAt(), now)
		eq(t, "Period.CompletedHours", back.Period().CompletedHours(), 4)
		eq(t, "Capacity.Max", back.Capacity().Max(), 30)
	})

	t.Run("courseSlot", func(t *testing.T) {
		s, _ := courseSlot.NewDayTime(16, 0)
		e, _ := courseSlot.NewDayTime(18, 30)
		tr, err := courseSlot.NewDayTimeRange(s, e)
		if err != nil {
			t.Fatal(err)
		}
		do := courseSlot.Reconstitute("slot-1", "c-1", time.Wednesday, tr, -1, "", now, now)

		po := conv.CourseSlotToPO(do)
		eq(t, "StartTime", po.StartTime, "16:00:00")
		eq(t, "EndTime", po.EndTime, "18:30:00")
		eq(t, "Weekday", po.Weekday, uint8(3))

		back, err := conv.CourseSlotToDO(po)
		if err != nil {
			t.Fatal(err)
		}
		eq(t, "DayTimeRange.Start", back.TimeRange().Start().String(), "16:00")
		eq(t, "DayTimeRange.End", back.TimeRange().End().String(), "18:30")
		eq(t, "CreatedAt", back.CreatedAt(), now)
	})

	t.Run("courseSlotChange", func(t *testing.T) {
		original := courseSlotChange.NewOriginalPlan(11, start, 1, "cl-1", "16:00", "18:00")
		target, err := courseSlotChange.NewTargetPlan(start, end, 2, "cl-2")
		if err != nil {
			t.Fatal(err)
		}
		do := courseSlotChange.Reconstitute(7, "c-1", 1, 2, original, target, "临时代课", now, now)

		po := conv.CourseSlotChangeToPO(do)
		eq(t, "OriginalStartTime", po.OriginalStartTime, "16:00")
		eq(t, "TargetTeacherID", po.TargetTeacherID, int64(2))

		back, err := conv.CourseSlotChangeToDO(po)
		if err != nil {
			t.Fatal(err)
		}
		eq(t, "OriginalPlan.Date", back.OriginalPlan().Date(), start)
		eq(t, "OriginalPlan.EndTimeStr", back.OriginalPlan().EndTimeStr(), "18:00")
		eq(t, "Reason", back.Reason(), "临时代课")
		eq(t, "UpdatedAt", back.UpdatedAt(), now)
	})

	t.Run("enrollment", func(t *testing.T) {
		do := enrollment.Reconstitute(1, 101, "c-1", 1, now, time.Time{}, time.Time{}, now)

		po := conv.EnrollmentToPO(do)
		eq(t, "CompletedAt.Valid", po.CompletedAt.Valid, false)

		back, err := conv.EnrollmentToDO(po)
		if err != nil {
			t.Fatal(err)
		}
		eq(t, "CompletedAt", back.CompletedAt(), time.Time{})
		eq(t, "EnrolledAt", back.EnrolledAt(), now)
		eq(t, "UpdatedAt", back.UpdatedAt(), now)
	})

	t.Run("makeup", func(t *testing.T) {
		do := makeup.Reconstitute(2, 101, "c-1", 11, start, 12, end, 2, 1, end, now, now)

		po := conv.MakeupToPO(do)
		eq(t, "CompletedAt.Valid", po.CompletedAt.Valid, true)

		back, err := conv.MakeupToDO(po)
		if err != nil {
			t.Fatal(err)
		}
		eq(t, "CompletedAt", back.CompletedAt(), end)
		eq(t, "OriginalSlotID", back.OriginalSlotID(), int64(11))
		eq(t, "TargetDate", back.TargetDate(), end)
	})

	t.Run("absence", func(t *testing.T) {
		do := absence.Reconstitute(1, 101, "c-1", 11, start, 2, 3, "家中有事", now, now)

		back, err := conv.AbsenceToDO(conv.AbsenceToPO(do))
		if err != nil {
			t.Fatal(err)
		}
		eq(t, "Reason", back.Reason(), "家中有事")
		eq(t, "AbsenceType", back.AbsenceType(), absence.AbsenceType(3))
		eq(t, "UpdatedAt", back.UpdatedAt(), now)
	})

	t.Run("qualification", func(t *testing.T) {
		do := qualification.Reconstitute(1, 1, "ct-1", now, 2, now)

		back, err := conv.QualificationToDO(conv.QualificationToPO(do))
		if err != nil {
			t.Fatal(err)
		}
		eq(t, "CourseTypeID", back.CourseTypeID(), "ct-1")
		eq(t, "UpdatedAt", back.UpdatedAt(), now)
	})
}

// TestDirtyDataRejected 读路径必须替领域模型挡住脏数据，而不是造出非法聚合。
func TestDirtyDataRejected(t *testing.T) {
	cases := []struct {
		name string
		run  func() error
	}{
		{"student: 手机号不合法", func() error {
			_, err := conv.StudentToDO(&model.Student{ID: 1, Name: "x", Phone: "12345", Status: 1})
			return err
		}},
		{"classroom: 楼栋为空", func() error {
			_, err := conv.ClassroomToDO(&model.Classroom{ID: "cl", Capacity: 30, Status: 1})
			return err
		}},
		{"course: 容量越界（enrolled > max）", func() error {
			_, err := conv.CourseToDO(&model.Course{
				ID: "c", CourseTypeID: "ct", CapacityMax: 10, CapacityEnrolled: 99,
				PeriodStartAt: now, PeriodEndAt: end, TotalHours: 40,
			})
			return err
		}},
		{"courseSlot: 时间段倒挂", func() error {
			_, err := conv.CourseSlotToDO(&model.CourseSlot{
				ID: "s", CourseID: "c", Weekday: 3,
				StartTime: "18:00:00", EndTime: "16:00:00",
			})
			return err
		}},
		{"courseSlotChange: 目标计划跨天", func() error {
			_, err := conv.CourseSlotChangeToDO(&model.CourseSlotChange{
				ID: 1, CourseID: "c", ApplicantID: 1,
				TargetStartAt: end, TargetEndAt: end.Add(24 * time.Hour),
				TargetTeacherID: 1, TargetClassroomID: "cl",
			})
			return err
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.run(); err == nil {
				t.Errorf("%s: 脏数据没有被拦下", tc.name)
			}
		})
	}
}
