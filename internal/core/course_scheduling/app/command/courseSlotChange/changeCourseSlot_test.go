package courseSlotChange

import (
	"context"
	"errors"
	"testing"
	"time"

	"cqrs/internal/core/course_scheduling/adapters/test_memory"
	commandmemory "cqrs/internal/core/course_scheduling/adapters/test_memory/command"
	memorycmd "cqrs/internal/core/course_scheduling/adapters/test_memory/command/implement"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
	"cqrs/internal/core/course_scheduling/domain/service/scheduleConflict"
)

// 测试用的资源与课程。目标时段统一取「下一个周一 09:00-11:00」。
const (
	targetTeacher   = int64(2)
	targetClassroom = "R102"
	ownCourse       = "C001"
	otherCourse     = "C999"
)

// newHandler 装配一个跑在干净内存存储上的命令处理器。
//
// 返回的是完整 store（测试要靠它塞数据和做断言），仓库拿到的则是收窄后的写侧面。
func newHandler(t *testing.T) (*Handler, *test_memory.Data) {
	t.Helper()

	store, cleanup, err := test_memory.NewData(nil)
	if err != nil {
		t.Fatalf("new data: %v", err)
	}
	t.Cleanup(cleanup)

	data := commandmemory.NewData(store)
	return NewHandler(
		memorycmd.NewCourseSlotChangeCommand(data),
		scheduleConflict.NewService(
			memorycmd.NewCourseSlotCommand(data),
			memorycmd.NewTeacherCommand(data),
			memorycmd.NewCourseCommand(data),
			memorycmd.NewCourseSlotChangeCommand(data),
		),
		data,
	), store
}

// nextMonday 返回下一个周一（严格晚于今天），避免目标时间落到过去。
func nextMonday() time.Time {
	d := time.Now().AddDate(0, 0, 1)
	for d.Weekday() != time.Monday {
		d = d.AddDate(0, 0, 1)
	}
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.Local)
}

// validCmd 造一条合法的临时换课命令：目标 = 下一个周一 09:00-11:00，讲师 2，教室 R102。
func validCmd(t *testing.T) ChangeCourseSlot {
	t.Helper()

	day := nextMonday()
	target, err := courseSlotChange.NewTargetPlan(
		day.Add(9*time.Hour), day.Add(11*time.Hour),
		targetTeacher, targetClassroom,
	)
	if err != nil {
		t.Fatalf("new target plan: %v", err)
	}

	return ChangeCourseSlot{
		CourseID:    ownCourse,
		ApplicantID: 1,
		ChangeType:  courseSlotChange.TypeSubstitute,
		Original: courseSlotChange.NewOriginalPlan(
			"550e8400-e29b-41d4-a716-446655440001", time.Now(), 1, "R101", "09:00", "11:00",
		),
		Target: target,
		Reason: "讲师出差，由李娜代课",
	}
}

// seedSlot 往内存存储里塞一条周排期。
func seedSlot(t *testing.T, d *test_memory.Data, courseID string, teacherID int64, classroomID string, weekday time.Weekday, fromHour, toHour int) {
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
	slot, err := courseSlot.NewCourseSlot(courseID, weekday, span, teacherID, classroomID)
	if err != nil {
		t.Fatalf("new course slot: %v", err)
	}

	d.SeedCourseSlot(slot)
}

// TestChangeCourseSlot 无冲突时登记成功，且落库一条。
func TestChangeCourseSlot(t *testing.T) {
	h, d := newHandler(t)
	cmd := validCmd(t)

	got, err := h.ChangeCourseSlot(context.Background(), cmd)
	if err != nil {
		t.Fatalf("ChangeCourseSlot: %v", err)
	}
	if got.CourseID() != cmd.CourseID || got.ApplicantID() != cmd.ApplicantID {
		t.Errorf("change = (%s, %d), want (%s, %d)",
			got.CourseID(), got.ApplicantID(), cmd.CourseID, cmd.ApplicantID)
	}
	if got.ChangeType() != courseSlotChange.TypeSubstitute {
		t.Errorf("changeType = %v, want %v", got.ChangeType(), courseSlotChange.TypeSubstitute)
	}
	if !got.IsSubstituteTeacher() {
		t.Error("IsSubstituteTeacher() = false, want true（原讲师 1 -> 目标讲师 2）")
	}
	if changes := d.CourseSlotChanges(); len(changes) != 1 {
		t.Errorf("变更单数量 = %d, want 1", len(changes))
	}
}

// TestChangeCourseSlotConflict 目标讲师或教室在目标时段被占用时应拒绝且不落库。
func TestChangeCourseSlotConflict(t *testing.T) {
	cases := []struct {
		name         string
		setup        func(t *testing.T, d *test_memory.Data)
		wantConflict bool
	}{
		{
			name: "目标讲师同期已有别的课",
			setup: func(t *testing.T, d *test_memory.Data) {
				// 讲师 2 周一 10:00-12:00 在别门课上 —— 与目标 09:00-11:00 重叠
				seedSlot(t, d, otherCourse, targetTeacher, "R201", time.Monday, 10, 12)
			},
			wantConflict: true,
		},
		{
			name: "目标教室同期已被别的课占用",
			setup: func(t *testing.T, d *test_memory.Data) {
				seedSlot(t, d, otherCourse, 7, targetClassroom, time.Monday, 10, 12)
			},
			wantConflict: true,
		},
		{
			name: "目标讲师同期有课但星期几不同",
			setup: func(t *testing.T, d *test_memory.Data) {
				seedSlot(t, d, otherCourse, targetTeacher, "R201", time.Tuesday, 9, 11)
			},
			wantConflict: false,
		},
		{
			name: "目标讲师同期有课但时间不重叠",
			setup: func(t *testing.T, d *test_memory.Data) {
				// 周一 13:00-15:00，与目标 09:00-11:00 不重叠
				seedSlot(t, d, otherCourse, targetTeacher, "R201", time.Monday, 13, 15)
			},
			wantConflict: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h, d := newHandler(t)
			tc.setup(t, d)

			got, err := h.ChangeCourseSlot(context.Background(), validCmd(t))

			if !tc.wantConflict {
				if err != nil {
					t.Fatalf("期望成功，实际 %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("期望被拒绝，实际成功: %+v", got)
			}
			if !errors.Is(err, courseSlotChange.ErrSlotChangeConflict) {
				t.Errorf("err = %v, want %v", err, courseSlotChange.ErrSlotChangeConflict)
			}
			if changes := d.CourseSlotChanges(); len(changes) != 0 {
				t.Errorf("被拒绝时不应写入，变更单数量 = %d", len(changes))
			}
		})
	}
}

// TestChangeCourseSlotIgnoresOwnCourseSlots 本门课自己的排期不算冲突 —— 那就是被换掉的那节课。
func TestChangeCourseSlotIgnoresOwnCourseSlots(t *testing.T) {
	h, d := newHandler(t)

	// 本门课自己的排期：同一讲师、同一教室、同一时段
	seedSlot(t, d, ownCourse, targetTeacher, targetClassroom, time.Monday, 9, 11)

	if _, err := h.ChangeCourseSlot(context.Background(), validCmd(t)); err != nil {
		t.Fatalf("本门课自己的排期不应判冲突，实际 %v", err)
	}
}

// TestChangeCourseSlotConflictWithOtherChange 别的换课记录占用了同一时段也应被拒。
func TestChangeCourseSlotConflictWithOtherChange(t *testing.T) {
	h, d := newHandler(t)
	day := nextMonday()

	// 已存在一条换课：讲师 2 周一 10:00-12:00 在 R202
	existing, err := courseSlotChange.NewTargetPlan(
		day.Add(10*time.Hour), day.Add(12*time.Hour), targetTeacher, "R202",
	)
	if err != nil {
		t.Fatalf("new target plan: %v", err)
	}
	d.SeedCourseSlotChange(courseSlotChange.Reconstitute(
		99, otherCourse, 1, courseSlotChange.TypeSubstitute,
		courseSlotChange.NewOriginalPlan("550e8400-e29b-41d4-a716-446655440004", day, 7, "R201", "10:00", "12:00"),
		existing, "讲师临时调整",
		time.Now(), time.Now(),
	))

	got, err := h.ChangeCourseSlot(context.Background(), validCmd(t))
	if err == nil {
		t.Fatalf("期望被拒绝，实际成功: %+v", got)
	}
	if !errors.Is(err, courseSlotChange.ErrSlotChangeConflict) {
		t.Errorf("err = %v, want %v", err, courseSlotChange.ErrSlotChangeConflict)
	}
	if changes := d.CourseSlotChanges(); len(changes) != 1 {
		t.Errorf("被拒绝时不应新增，变更单数量 = %d, want 1", len(changes))
	}
}

// TestChangeCourseSlotRejected 聚合校验不过时不落库。
func TestChangeCourseSlotRejected(t *testing.T) {
	// 固定用「两天前的 09:00-11:00」，别拿 time.Now() 直接加小时：
	// 22:00 之后 +2h 会翻到第二天，聚合构造函数会以「起止必须同一天」拒掉，
	// 于是这条用例每到深夜就红。
	twoDaysAgo := time.Now().AddDate(0, 0, -2)
	day := time.Date(twoDaysAgo.Year(), twoDaysAgo.Month(), twoDaysAgo.Day(), 0, 0, 0, 0, time.Local)
	past, err := courseSlotChange.NewTargetPlan(
		day.Add(9*time.Hour),
		day.Add(11*time.Hour),
		2, "R102",
	)
	if err != nil {
		t.Fatalf("new target plan: %v", err)
	}

	cases := []struct {
		name    string
		mutate  func(cmd *ChangeCourseSlot)
		wantErr error
	}{
		{"课程为空", func(cmd *ChangeCourseSlot) { cmd.CourseID = "" }, nil},
		{"申请人为 0", func(cmd *ChangeCourseSlot) { cmd.ApplicantID = 0 }, nil},
		{"事由为空", func(cmd *ChangeCourseSlot) { cmd.Reason = "" }, nil},
		{"目标时间在过去", func(cmd *ChangeCourseSlot) { cmd.Target = past }, courseSlotChange.ErrTargetDateInPast},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h, d := newHandler(t)
			cmd := validCmd(t)
			tc.mutate(&cmd)

			got, err := h.ChangeCourseSlot(context.Background(), cmd)
			if err == nil {
				t.Fatalf("期望被拒绝，实际成功: %+v", got)
			}
			if tc.wantErr != nil && !errors.Is(err, tc.wantErr) {
				t.Errorf("err = %v, want %v", err, tc.wantErr)
			}
			if changes := d.CourseSlotChanges(); len(changes) != 0 {
				t.Errorf("被拒绝时不应写入，变更单数量 = %d", len(changes))
			}
		})
	}
}

// TestDeleteCourseSlotChange 删除不存在的变更单返回 not found。
func TestDeleteCourseSlotChange(t *testing.T) {
	h, _ := newHandler(t)

	got, err := h.ChangeCourseSlot(context.Background(), validCmd(t))
	if err != nil {
		t.Fatalf("ChangeCourseSlot: %v", err)
	}

	if err := h.ChangeCmd.DeleteCourseSlotChange(context.Background(), got.ID()); err != nil {
		t.Errorf("Delete(existing) = %v, want nil", err)
	}
	if err := h.ChangeCmd.DeleteCourseSlotChange(context.Background(), got.ID()); !errors.Is(err, courseSlotChange.ErrSlotChangeNotFound) {
		t.Errorf("Delete(missing) = %v, want %v", err, courseSlotChange.ErrSlotChangeNotFound)
	}
}
