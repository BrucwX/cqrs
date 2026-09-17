package enrollment

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"cqrs/internal/core/course_scheduling/adapters/memImp4test"
	commandmemory "cqrs/internal/core/course_scheduling/adapters/memImp4test/command"
	memorycmd "cqrs/internal/core/course_scheduling/adapters/memImp4test/command/implement"
	querymemory "cqrs/internal/core/course_scheduling/adapters/memImp4test/query"
	memoryquery "cqrs/internal/core/course_scheduling/adapters/memImp4test/query/implement"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
	"cqrs/internal/core/course_scheduling/domain/service/scheduleConflict"
)

const (
	// courseFree 窗口内、未满、周二 09:00-11:00，与学员 A 现有课程不撞；
	// 学员 A 有一条「未选课」注册记录 → 可以正式选课
	courseFree = "c-free"
	// courseTaken 学员 A 已在学，周一 09:00-11:00
	courseTaken = "c-taken"
	// courseClash 窗口内、未满，周一 10:00-12:00（与 courseTaken 时间重叠）
	courseClash = "c-clash"
	// courseFull 窗口内，但已满
	courseFull = "c-full"
	// courseClosed 窗口已关闭
	courseClosed = "c-closed"
	// courseRetake 窗口内，学员 A 曾选过但已退课，且已重新登记一条未选课记录 → 可以重选
	courseRetake = "c-retake"
	// courseDupe 窗口内，学员 A 已在学（同一门课重复选 → 走时间冲突规则拦下）
	courseDupe = "c-dupe"
	// courseNoQuota 窗口内、未满、不撞时间，但学员 A 一条注册记录都没有 → 没资格
	courseNoQuota = "c-no-quota"
	// courseRepay 窗口内、未满、不撞时间，但学员 A 只有一条「已退课」记录（没重新登记）→ 没资格
	courseRepay = "c-repay"

	studentA = int64(1)
)

// newData 造一份干净的内存数据：
// 9 门课 + 9 条排期 + 学员 A 的 11 条注册记录。
//
// 注册记录分三类：在学（courseTaken / courseDupe）、已退课（courseRetake 的旧记录、courseRepay）、
// 未选课（= 已报班/已缴费的选课资格）。选课要求先有未选课记录，所以 courseFree /
// courseClash / courseFull / courseClosed / courseDupe / courseRetake / c-not-exist 各挂了一条。
func newData(t *testing.T) *memImp4test.Data {
	t.Helper()

	d, cleanup, err := memImp4test.NewData(nil)
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
		newCourse(t, courseNoQuota, openWindow, 10, 0),
		newCourse(t, courseRepay, openWindow, 10, 0),
	)
	d.SeedCourseSlot(
		newSlot(t, courseFree, time.Tuesday, 9, 11),
		newSlot(t, courseTaken, time.Monday, 9, 11),
		newSlot(t, courseClash, time.Monday, 10, 12),
		newSlot(t, courseFull, time.Friday, 9, 11),
		newSlot(t, courseClosed, time.Saturday, 9, 11),
		newSlot(t, courseRetake, time.Thursday, 9, 11),
		newSlot(t, courseDupe, time.Wednesday, 9, 11),
		newSlot(t, courseNoQuota, time.Tuesday, 13, 15),
		newSlot(t, courseRepay, time.Thursday, 13, 15),
	)

	// 学员 A：在学 2 条、已退课 2 条、未选课 7 条
	d.SeedEnrollment(
		enrollment.Reconstitute(
			1, studentA, courseTaken, enrollment.StatusEnrolled,
			now, time.Time{}, time.Time{}, now,
		),
		enrollment.Reconstitute(
			2, studentA, courseDupe, enrollment.StatusEnrolled,
			now, time.Time{}, time.Time{}, now,
		),
		// courseRetake：先退课，之后重新登记了一条资格（ID 4）
		enrollment.Reconstitute(
			3, studentA, courseRetake, enrollment.StatusDropped,
			now, time.Time{}, now, now,
		),
		enrollment.Reconstitute(
			4, studentA, courseRetake, enrollment.StatusNotSelected,
			now, time.Time{}, time.Time{}, now,
		),
		// courseRepay：只有已退课，没重新登记 → 不能重选
		enrollment.Reconstitute(
			5, studentA, courseRepay, enrollment.StatusDropped,
			now, time.Time{}, now, now,
		),
		// 未选课（选课资格）
		enrollment.Reconstitute(
			6, studentA, courseFree, enrollment.StatusNotSelected,
			now, time.Time{}, time.Time{}, now,
		),
		enrollment.Reconstitute(
			7, studentA, courseClash, enrollment.StatusNotSelected,
			now, time.Time{}, time.Time{}, now,
		),
		enrollment.Reconstitute(
			8, studentA, courseFull, enrollment.StatusNotSelected,
			now, time.Time{}, time.Time{}, now,
		),
		enrollment.Reconstitute(
			9, studentA, courseClosed, enrollment.StatusNotSelected,
			now, time.Time{}, time.Time{}, now,
		),
		enrollment.Reconstitute(
			10, studentA, courseDupe, enrollment.StatusNotSelected,
			now, time.Time{}, time.Time{}, now,
		),
		enrollment.Reconstitute(
			11, studentA, "c-not-exist", enrollment.StatusNotSelected,
			now, time.Time{}, time.Time{}, now,
		),
	)

	return d
}

// TestStudentEnroll_Allowed 有资格且三项检查都通过时，把那条未选课记录推到在读（不新建记录）。
func TestStudentEnroll_Allowed(t *testing.T) {
	cases := []struct {
		name     string
		courseID string
		// wantID 应当被流转的那条未选课记录
		wantID int64
		// wantKept 不该被碰的那条记录（0 = 没有）
		wantKept int64
	}{
		{"已有未选课记录", courseFree, 6, 0},
		{"退课后重新登记再选", courseRetake, 4, 3},
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
			if got.ID() != tc.wantID {
				t.Errorf("流转的记录 ID = %d, want %d", got.ID(), tc.wantID)
			}
			if got.CourseID() != tc.courseID || got.StudentID() != studentA {
				t.Errorf("enrollment = (%d, %s), want (%d, %s)",
					got.StudentID(), got.CourseID(), studentA, tc.courseID)
			}
			if !got.IsActive() {
				t.Errorf("选课后应当是 StatusEnrolled，实际 %v", got.Status())
			}
			// 写回的是库里那条，不是只改了内存副本
			if stored := findByID(t, d, tc.wantID); !stored.IsActive() {
				t.Errorf("库里的记录 %d 状态 = %v, want %v",
					tc.wantID, stored.Status(), enrollment.StatusEnrolled)
			}
			// 不再新建记录
			if after := len(d.Enrollments()); after != before {
				t.Errorf("选课不应新建记录，记录数 = %d, want %d", after, before)
			}
			// 旧的已退课记录不受影响
			if tc.wantKept != 0 {
				if stored := findByID(t, d, tc.wantKept); stored.Status() != enrollment.StatusDropped {
					t.Errorf("记录 %d 不应被改动，状态 = %v", tc.wantKept, stored.Status())
				}
			}
		})
	}
}

// TestStudentEnroll_Rejected 没资格 或 任一规则不过都应拒绝且不落库。
func TestStudentEnroll_Rejected(t *testing.T) {
	cases := []struct {
		name     string
		courseID string
		wantErr  error
	}{
		// 没有「未选课」的注册记录 = 这门课还没缴费
		{"一条记录都没有", courseNoQuota, enrollment.ErrEnrollmentNotPaid},
		{"只有已退课记录（没重新缴费）", courseRepay, enrollment.ErrEnrollmentNotPaid},
		// 有资格但准入不过
		{"重复选课", courseDupe, enrollment.ErrEnrollmentConflict},
		{"窗口已关闭", courseClosed, enrollment.ErrEnrollmentConflict},
		{"课程已满", courseFull, enrollment.ErrEnrollmentConflict},
		{"与在学课程撞时间", courseClash, enrollment.ErrEnrollmentConflict},
		// 课程取不到时仓库会直接报 not found，不再当成冲突
		{"课程不存在", "c-not-exist", course.ErrCourseNotFound},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := newData(t)
			h := newHandler(d)

			before := len(d.Enrollments())
			beforeStatus := fmt.Sprint(enrollmentStatuses(d, studentA, tc.courseID))

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
			if afterStatus := fmt.Sprint(enrollmentStatuses(d, studentA, tc.courseID)); afterStatus != beforeStatus {
				t.Errorf("被拒绝时状态不应变化：%s -> %s", beforeStatus, afterStatus)
			}
		})
	}
}

// --- 测试辅助 ---

// findByID 按 ID 取库里的注册记录（取不到直接失败）。
func findByID(t *testing.T, d *memImp4test.Data, id int64) enrollment.CourseEnrollment {
	t.Helper()

	for _, e := range d.Enrollments() {
		if e.ID() == id {
			return *e
		}
	}

	t.Fatalf("注册记录 %d 不存在", id)
	return enrollment.CourseEnrollment{}
}

// enrollmentStatuses 该学员在该课程下的全部记录状态（按 ID 升序，用于比对「没被改动」）。
func enrollmentStatuses(d *memImp4test.Data, studentID int64, courseID string) []string {
	out := make([]string, 0, 1)
	for _, e := range d.Enrollments() {
		if e.StudentID() == studentID && e.CourseID() == courseID {
			out = append(out, fmt.Sprintf("%d:%v", e.ID(), e.Status()))
		}
	}
	return out
}

// newHandler 装配一个处理器：写侧仓库拿写侧窄面，读侧仓库拿读侧窄面。
//
// 这个用例两侧都用（报要和查在读侧），所以两个窄面各收一次。
func newHandler(d *memImp4test.Data) *Handler {
	cmdData := commandmemory.NewData(d)
	queryData := querymemory.NewData(d)

	return NewHandler(
		memorycmd.NewEnrollmentCommand(cmdData),
		memoryquery.NewCourseEnrollmentQuery(queryData),
		memoryquery.NewStudentQuery(queryData),
		memoryquery.NewCourseQuery(queryData),
		scheduleConflict.NewService(
			memorycmd.NewCourseSlotCommand(cmdData),
			memorycmd.NewTeacherCommand(cmdData),
			memorycmd.NewCourseCommand(cmdData),
			memorycmd.NewCourseSlotChangeCommand(cmdData),
		),
		cmdData,
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
