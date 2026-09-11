package enrollment

import (
	"context"
	"errors"
	"testing"
	"time"

	"cqrs/internal/core/course_scheduling/adapters/memory"
	memorycmd "cqrs/internal/core/course_scheduling/adapters/memory/command"
	memoryquery "cqrs/internal/core/course_scheduling/adapters/memory/query"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

const (
	// courseFree 窗口内、未满、周二 09:00-11:00，与学员 A 现有课程不撞
	courseFree = "c-free"
	// courseTaken 学员 A 已在学，周一 09:00-11:00
	courseTaken = "c-taken"
	// courseClash 窗口内、未满，周一 10:00-12:00（与 courseTaken 时间重叠）
	courseClash = "c-clash"
	// courseFull 窗口内，但已满
	courseFull = "c-full"
	// courseClosed 窗口已关闭
	courseClosed = "c-closed"
	// courseRetake 窗口内，学员 A 曾选过但已退课 → 可以重选
	courseRetake = "c-retake"
	// courseDupe 窗口内，学员 A 已在学（同一门课重复选 → 走时间冲突规则拦下）
	courseDupe = "c-dupe"

	studentA = int64(1)
)

// newData 造一份干净的内存数据：
// 7 门课 + 7 条排期 + 学员 A 的 3 条注册记录（courseTaken/courseDupe 在学，courseRetake 已退课）。
func newData(t *testing.T) *memory.Data {
	t.Helper()

	d, cleanup, err := memory.NewData(nil)
	if err != nil {
		t.Fatalf("new data: %v", err)
	}
	t.Cleanup(cleanup)

	now := time.Now()
	openWindow := course.NewEnrollmentWindow(
		now.AddDate(0, 0, -1), now.AddDate(0, 0, 7), now.AddDate(0, 0, 14),
	)
	closedWindow := course.NewEnrollmentWindow(
		now.AddDate(0, 0, -30), now.AddDate(0, 0, -10), now.AddDate(0, 0, -5),
	)

	d.SeedCourse(
		newCourse(t, courseFree, openWindow, 10, 0),
		newCourse(t, courseTaken, openWindow, 10, 1),
		newCourse(t, courseClash, openWindow, 10, 0),
		newCourse(t, courseFull, openWindow, 1, 1),
		newCourse(t, courseClosed, closedWindow, 10, 0),
		newCourse(t, courseRetake, openWindow, 10, 0),
		newCourse(t, courseDupe, openWindow, 10, 1),
	)
	d.SeedCourseSlot(
		newSlot(t, courseFree, time.Tuesday, 9, 11),
		newSlot(t, courseTaken, time.Monday, 9, 11),
		newSlot(t, courseClash, time.Monday, 10, 12),
		newSlot(t, courseFull, time.Friday, 9, 11),
		newSlot(t, courseClosed, time.Saturday, 9, 11),
		newSlot(t, courseRetake, time.Thursday, 9, 11),
		newSlot(t, courseDupe, time.Wednesday, 9, 11),
	)

	// 学员 A：courseTaken、courseDupe 在学；courseRetake 已退课
	d.SeedEnrollment(
		enrollment.Reconstitute(
			1, studentA, courseTaken, enrollment.StatusEnrolled,
			now, time.Time{}, time.Time{}, now,
		),
		enrollment.Reconstitute(
			2, studentA, courseDupe, enrollment.StatusEnrolled,
			now, time.Time{}, time.Time{}, now,
		),
		enrollment.Reconstitute(
			3, studentA, courseRetake, enrollment.StatusDropped,
			now, time.Time{}, now, now,
		),
	)

	return d
}

// TestStudentEnroll_Allowed 四项检查都通过时应当写入。
func TestStudentEnroll_Allowed(t *testing.T) {
	cases := []struct {
		name     string
		courseID string
	}{
		{"首次选课", courseFree},
		{"退课后重新选课", courseRetake},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := newData(t)
			h := newHandler(d)

			before := len(d.Enrollments())

			got, err := h.StudentEnroll(context.Background(), StudentEnroll{
				StudentID: studentA,
				CourseID:  tc.courseID,
			})
			if err != nil {
				t.Fatalf("StudentEnroll: %v", err)
			}
			if got.CourseID() != tc.courseID || got.StudentID() != studentA {
				t.Errorf("enrollment = (%d, %s), want (%d, %s)",
					got.StudentID(), got.CourseID(), studentA, tc.courseID)
			}
			if !got.IsActive() {
				t.Errorf("新选课应当是 StatusEnrolled，实际 %v", got.Status())
			}
			if after := len(d.Enrollments()); after != before+1 {
				t.Errorf("选课记录数 = %d, want %d", after, before+1)
			}
		})
	}
}

// TestStudentEnroll_Rejected 三条规则任一不过都应拒绝且不落库。
func TestStudentEnroll_Rejected(t *testing.T) {
	cases := []struct {
		name     string
		courseID string
		wantErr  error
	}{
		{"重复选课", courseDupe, repo.ErrEnrollmentConflict},
		{"窗口已关闭", courseClosed, repo.ErrEnrollmentConflict},
		{"课程已满", courseFull, repo.ErrEnrollmentConflict},
		{"与在学课程撞时间", courseClash, repo.ErrEnrollmentConflict},
		// 课程取不到时仓库会直接报 not found，不再当成冲突
		{"课程不存在", "c-not-exist", repo.ErrCourseNotFound},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := newData(t)
			h := newHandler(d)

			before := len(d.Enrollments())

			got, err := h.StudentEnroll(context.Background(), StudentEnroll{
				StudentID: studentA,
				CourseID:  tc.courseID,
			})
			if err == nil {
				t.Fatalf("期望被拒绝，实际选课成功: %+v", got)
			}
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("err = %v, want 包含 %v", err, tc.wantErr)
			}
			if after := len(d.Enrollments()); after != before {
				t.Errorf("被拒绝时不应写入，记录数 = %d, want %d", after, before)
			}
		})
	}
}

// --- 测试辅助 ---

func newHandler(d *memory.Data) *Handler {
	return NewHandler(
		memorycmd.NewEnrollmentCommand(d),
		memoryquery.NewCourseEnrollmentQuery(d),
		memoryquery.NewStudentQuery(d),
		memoryquery.NewCourseQuery(d),
		memorycmd.NewEnrollRepo(d),
	)
}

func newCourse(
	t *testing.T,
	id string,
	window course.EnrollmentWindow,
	max int,
	enrolled int,
) *course.Course {
	t.Helper()

	capacity, err := course.NewCapacity(max, enrolled)
	if err != nil {
		t.Fatalf("new capacity: %v", err)
	}

	now := time.Now()
	period, err := course.NewCoursePeriod(
		now.AddDate(0, 0, -7), now.AddDate(0, 6, 0), 16, 0,
	)
	if err != nil {
		t.Fatalf("new period: %v", err)
	}

	return course.Reconstitute(id, "ct-demo", capacity, window, period)
}

func newSlot(
	t *testing.T,
	courseID string,
	weekday time.Weekday,
	fromHour int,
	toHour int,
) *courseSlot.CourseSlot {
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

	slot, err := courseSlot.NewCourseSlot(courseID, weekday, span, -1, "")
	if err != nil {
		t.Fatalf("new course slot: %v", err)
	}

	return slot
}
