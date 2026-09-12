package qualification

import (
	"context"
	"errors"
	"testing"
	"time"

	"cqrs/internal/core/course_scheduling/adapters/memorystore"
	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
)

const (
	teacherWang = int64(3)   // 讲师 ID
	studentWang = int64(203) // 王强作为学员上课时的 ID

	ctProgramming  = "ct-0001" // 少儿编程
	ctEnglish      = "ct-0002" // 成人英语
	courseProgram  = "C001"    // 属于 ct-0001
	courseProgram2 = "C002"    // 也属于 ct-0001
	courseEnglish  = "C003"    // 属于 ct-0002
)

// seedTeacher 造一个讲师，并给他一个「作为学员上课」的 ID。
func seedTeacher(t *testing.T, d *memorystore.Data, id, studentID int64) {
	t.Helper()

	now := time.Now()
	d.SeedTeacher(teacher.Reconstitute(
		id, studentID, "王强", "讲师", teacher.ContactInfo{}, teacher.StatusActive, now, now,
	))
}

// seedCourse 造一门属于某课程类型的课程（授证判定只关心课程 ID 与类型）。
func seedCourse(t *testing.T, d *memorystore.Data, id, courseTypeID string) {
	t.Helper()

	now := time.Now()
	capacity, err := course.NewCapacity(10, 0)
	if err != nil {
		t.Fatalf("new capacity: %v", err)
	}
	period, err := course.NewCoursePeriod(now.AddDate(0, -2, 0), now, 20, 20)
	if err != nil {
		t.Fatalf("new period: %v", err)
	}

	d.SeedCourse(course.Reconstitute(
		id,
		courseTypeID,
		capacity,
		course.NewEnrollmentWindow(now.AddDate(0, -4, 0), now.AddDate(0, -3, 0), now.AddDate(0, -3, 0)),
		period,
	))
}

// seedEnrollment 造一条该讲师（以学员身份参训）的报名记录。
func seedEnrollment(t *testing.T, d *memorystore.Data, id int64, courseID string, status enrollment.Status) {
	t.Helper()

	now := time.Now()
	d.SeedEnrollment(enrollment.Reconstitute(
		id, studentWang, courseID, status, now, now, time.Time{}, now,
	))
}

// seedAbsence 造一条该讲师在某门课上的缺勤记录。
func seedAbsence(t *testing.T, d *memorystore.Data, id int64, courseID string) {
	t.Helper()

	now := time.Now()
	d.SeedAbsence(absence.Reconstitute(
		id, studentWang, courseID, "550e8400-e29b-41d4-a716-446655440002", now, 2,
		absence.TypePersonalLeave, "事假", now, now,
	))
}

// TestQualifyGrants 修完该类型下的课且没缺勤 -> 发证。
func TestQualifyGrants(t *testing.T) {
	h, d := newHandler(t)
	seedTeacher(t, d, teacherWang, studentWang)
	seedCourse(t, d, courseProgram, ctProgramming)
	seedEnrollment(t, d, 1, courseProgram, enrollment.StatusCompleted)

	if err := h.Qualify(context.Background(), QualifyInput{
		TeacherID:    teacherWang,
		CourseTypeID: ctProgramming,
	}); err != nil {
		t.Fatalf("Qualify: %v", err)
	}

	items := d.Qualifications()
	if len(items) != 1 {
		t.Fatalf("资质数量 = %d, want 1", len(items))
	}
	got := items[0]
	if got.ID() == 0 {
		t.Error("资质应当自动分配 ID")
	}
	if got.TeacherID() != teacherWang || got.CourseTypeID() != ctProgramming {
		t.Errorf("qualification = (%d, %q), want (%d, %q)",
			got.TeacherID(), got.CourseTypeID(), teacherWang, ctProgramming)
	}
	if got.Status() != qualification.StatusActive {
		t.Errorf("status = %v, want %v", got.Status(), qualification.StatusActive)
	}
	if got.CertifiedAt().IsZero() {
		t.Error("certifiedAt 应当取发证时间")
	}
}

// TestQualifyGrantsBySiblingCourse 该类型下另一门课修完也算（缺勤只看那一门）。
func TestQualifyGrantsBySiblingCourse(t *testing.T) {
	h, d := newHandler(t)
	seedTeacher(t, d, teacherWang, studentWang)
	seedCourse(t, d, courseProgram, ctProgramming)
	seedCourse(t, d, courseProgram2, ctProgramming)
	seedEnrollment(t, d, 1, courseProgram, enrollment.StatusCompleted)
	seedEnrollment(t, d, 2, courseProgram2, enrollment.StatusEnrolled)
	seedAbsence(t, d, 1, courseProgram2) // 缺的是另一门课

	if err := h.Qualify(context.Background(), QualifyInput{
		TeacherID:    teacherWang,
		CourseTypeID: ctProgramming,
	}); err != nil {
		t.Fatalf("Qualify: %v", err)
	}

	if got := len(d.Qualifications()); got != 1 {
		t.Errorf("资质数量 = %d, want 1", got)
	}
}

// TestQualifyRejectsUnfinished 没结业 -> 拒绝，且不落库。
func TestQualifyRejectsUnfinished(t *testing.T) {
	cases := map[string]enrollment.Status{
		"在读":  enrollment.StatusEnrolled,
		"已退课": enrollment.StatusDropped,
	}

	for name, status := range cases {
		t.Run(name, func(t *testing.T) {
			h, d := newHandler(t)
			seedTeacher(t, d, teacherWang, studentWang)
			seedCourse(t, d, courseProgram, ctProgramming)
			seedEnrollment(t, d, 1, courseProgram, status)

			err := h.Qualify(context.Background(), QualifyInput{
				TeacherID:    teacherWang,
				CourseTypeID: ctProgramming,
			})
			if !errors.Is(err, qualification.ErrCourseNotFinished) {
				t.Errorf("err = %v, want %v", err, qualification.ErrCourseNotFinished)
			}
			if got := len(d.Qualifications()); got != 0 {
				t.Errorf("资质数量 = %d, want 0", got)
			}
		})
	}
}

// TestQualifyRejectsAbsence 结业了但那门课缺过勤 -> 拒绝。
func TestQualifyRejectsAbsence(t *testing.T) {
	h, d := newHandler(t)
	seedTeacher(t, d, teacherWang, studentWang)
	seedCourse(t, d, courseProgram, ctProgramming)
	seedEnrollment(t, d, 1, courseProgram, enrollment.StatusCompleted)
	seedAbsence(t, d, 1, courseProgram)

	err := h.Qualify(context.Background(), QualifyInput{
		TeacherID:    teacherWang,
		CourseTypeID: ctProgramming,
	})
	if !errors.Is(err, qualification.ErrCourseNotFinished) {
		t.Errorf("err = %v, want %v", err, qualification.ErrCourseNotFinished)
	}
	if got := len(d.Qualifications()); got != 0 {
		t.Errorf("资质数量 = %d, want 0", got)
	}
}

// TestQualifyRejectsOtherCourseType 修完的是别的类型的课 -> 拒绝。
func TestQualifyRejectsOtherCourseType(t *testing.T) {
	h, d := newHandler(t)
	seedTeacher(t, d, teacherWang, studentWang)
	seedCourse(t, d, courseEnglish, ctEnglish)
	seedEnrollment(t, d, 1, courseEnglish, enrollment.StatusCompleted)

	err := h.Qualify(context.Background(), QualifyInput{
		TeacherID:    teacherWang,
		CourseTypeID: ctProgramming,
	})
	if !errors.Is(err, qualification.ErrCourseNotFinished) {
		t.Errorf("err = %v, want %v", err, qualification.ErrCourseNotFinished)
	}
}

// TestQualifyRejectsNoCourse 该类型下压根没有课程 -> 拒绝。
func TestQualifyRejectsNoCourse(t *testing.T) {
	h, d := newHandler(t)
	seedTeacher(t, d, teacherWang, studentWang)

	err := h.Qualify(context.Background(), QualifyInput{
		TeacherID:    teacherWang,
		CourseTypeID: ctProgramming,
	})
	if !errors.Is(err, qualification.ErrCourseNotFinished) {
		t.Errorf("err = %v, want %v", err, qualification.ErrCourseNotFinished)
	}
}

// TestQualifyRejectsUnknownTeacher 讲师不存在 -> 连他的学员 ID 都取不到。
func TestQualifyRejectsUnknownTeacher(t *testing.T) {
	h, d := newHandler(t)
	seedCourse(t, d, courseProgram, ctProgramming)
	seedEnrollment(t, d, 1, courseProgram, enrollment.StatusCompleted)

	err := h.Qualify(context.Background(), QualifyInput{
		TeacherID:    teacherWang,
		CourseTypeID: ctProgramming,
	})
	if !errors.Is(err, teacher.ErrTeacherNotFound) {
		t.Errorf("err = %v, want %v", err, teacher.ErrTeacherNotFound)
	}
}

// TestQualifyRejectsOtherStudent 报名/缺勤记录挂的是学员 ID，不是讲师 ID。
func TestQualifyRejectsOtherStudent(t *testing.T) {
	h, d := newHandler(t)
	seedTeacher(t, d, teacherWang, studentWang)
	seedCourse(t, d, courseProgram, ctProgramming)

	now := time.Now()
	d.SeedEnrollment(enrollment.Reconstitute(
		1, teacherWang, courseProgram, enrollment.StatusCompleted, now, now, time.Time{}, now,
	))

	err := h.Qualify(context.Background(), QualifyInput{
		TeacherID:    teacherWang,
		CourseTypeID: ctProgramming,
	})
	if !errors.Is(err, qualification.ErrCourseNotFinished) {
		t.Errorf("err = %v, want %v", err, qualification.ErrCourseNotFinished)
	}
}

// TestQualifyRejectsInvalidInput 讲师 ID / 课程类型 ID 缺失 -> 报参数错误。
func TestQualifyRejectsInvalidInput(t *testing.T) {
	cases := map[string]struct {
		cmd  QualifyInput
		want error
	}{
		"没有讲师":   {QualifyInput{CourseTypeID: ctProgramming}, qualification.ErrTeacherRequired},
		"没有课程类型": {QualifyInput{TeacherID: teacherWang}, qualification.ErrCourseTypeRequired},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			h, d := newHandler(t)
			seedTeacher(t, d, teacherWang, studentWang)
			seedCourse(t, d, courseProgram, ctProgramming)
			seedEnrollment(t, d, 1, courseProgram, enrollment.StatusCompleted)

			if err := h.Qualify(context.Background(), tc.cmd); !errors.Is(err, tc.want) {
				t.Errorf("err = %v, want %v", err, tc.want)
			}
			if got := len(d.Qualifications()); got != 0 {
				t.Errorf("资质数量 = %d, want 0", got)
			}
		})
	}
}
