package makeup

import (
	"context"
	"errors"
	"testing"
	"time"

	commandmemory "cqrs/internal/core/course_scheduling/adapters/command/memory"
	memorycmd "cqrs/internal/core/course_scheduling/adapters/command/memory/implement"
	"cqrs/internal/core/course_scheduling/adapters/memorystore"
	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
	"cqrs/internal/core/course_scheduling/domain/service/classroomCapacity"
)

const (
	// 补的那门课 20 人；room 30 座装得下（还得给蹭课的留 1 个位），smallRoom 5 座装不下。
	makeupCourse  = "C001"
	courseSeats   = 20
	room          = "R101"
	roomSeats     = 30
	smallRoom     = "R102"
	smallSeats    = 5
	makeupTeacher = int64(1)
	makeupStudent = int64(101)
)

// fixture 是一份搭好的内存数据 + 处理器。
//
// 两个目标槽位属于同一门课、排在两间教室：targetSlot 那间装得下，smallSlot 装不下。
type fixture struct {
	handler    *Handler
	data       *memorystore.Data
	targetSlot string // 排在 30 座的教室
	smallSlot  string // 排在 5 座的小教室
}

// newFixture 造一份数据：
// 1 门 20 人的课、2 间教室（30 座 / 5 座）、2 条排期（各排一间）。
func newFixture(t *testing.T) *fixture {
	t.Helper()

	store, cleanup, err := memorystore.NewData(nil)
	if err != nil {
		t.Fatalf("new data: %v", err)
	}
	t.Cleanup(cleanup)

	seedCourse(t, store, makeupCourse, courseSeats)
	seedClassroom(t, store, room, roomSeats)
	seedClassroom(t, store, smallRoom, smallSeats)

	target := seedSlot(t, store, makeupCourse, time.Monday, 9, 11, room)
	small := seedSlot(t, store, makeupCourse, time.Tuesday, 9, 11, smallRoom)

	data := commandmemory.NewData(store)
	makeupCmd := memorycmd.NewMakeupCommand(data)
	return &fixture{
		handler: NewHandler(
			makeupCmd,
			classroomCapacity.NewService(
				memorycmd.NewClassroomCommand(data),
				memorycmd.NewCourseSlotCommand(data),
				memorycmd.NewCourseCommand(data),
				makeupCmd,
			),
			data,
		),
		data:       store,
		targetSlot: target,
		smallSlot:  small,
	}
}

// validCmd 造一条能过的补课命令：补到 targetSlot 那节课。
func (f *fixture) validCmd() RecordMakeup {
	return RecordMakeup{
		StudentID:      makeupStudent,
		CourseID:       makeupCourse,
		OriginalSlotID: "550e8400-e29b-41d4-a716-446655440099",
		OriginalDate:   time.Date(2026, time.September, 14, 0, 0, 0, 0, time.Local),
		TargetSlotID:   f.targetSlot,
		TargetDate:     time.Date(2026, time.September, 16, 0, 0, 0, 0, time.Local),
		MakeupHours:    2,
	}
}

// seedMakeups 往 (slotID, date) 这节课上塞 n 条别的学员的补课记录。
//
// 状态由调用方给 —— 判定只该认 StatusBooked 的那些。
func (f *fixture) seedMakeups(t *testing.T, n int, slotID string, date time.Time, status makeup.Status) {
	t.Helper()

	for i := range n {
		f.data.SeedMakeup(makeup.Reconstitute(
			int64(1000+i), int64(200+i), makeupCourse,
			"550e8400-e29b-41d4-a716-446655440099", date.AddDate(0, 0, -1),
			slotID, date, 2,
			status, time.Time{}, time.Now(), time.Now(),
		))
	}
}

// --- 造数据的辅助 ---

// seedCourse 塞一门课进去（容量判定只用得到 ID 与人数上限）。
func seedCourse(t *testing.T, d *memorystore.Data, id string, max int) {
	t.Helper()

	capacity, err := course.NewCapacity(max, 0)
	if err != nil {
		t.Fatalf("new capacity: %v", err)
	}
	now := time.Now()
	period, err := course.NewCoursePeriod(now.AddDate(0, -1, 0), now.AddDate(0, 2, 0), 16, 0)
	if err != nil {
		t.Fatalf("new period: %v", err)
	}

	d.SeedCourse(course.Reconstitute(
		id, "ct-demo", capacity,
		course.NewEnrollmentWindow(now.AddDate(0, -2, 0), now.AddDate(0, 2, 0), now.AddDate(0, 3, 0)),
		period,
	))
}

// seedClassroom 塞一间教室进去；ID 固定所以走 Reconstitute。
func seedClassroom(t *testing.T, d *memorystore.Data, id string, seats int) {
	t.Helper()

	loc, err := classroom.NewLocation("A", 1, id)
	if err != nil {
		t.Fatalf("new location: %v", err)
	}
	d.SeedClassroom(classroom.Reconstitute(id, loc, seats, 0, classroom.StatusAvailable))
}

// seedSlot 塞一条指定教室的排期进去，返回它的 ID（聚合自己生成的 UUID）。
func seedSlot(t *testing.T, d *memorystore.Data, courseID string, weekday time.Weekday, fromHour, toHour int, classroomID string) string {
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
	slot, err := courseSlot.NewCourseSlot(courseID, weekday, span, makeupTeacher, classroomID)
	if err != nil {
		t.Fatalf("new course slot: %v", err)
	}

	d.SeedCourseSlot(slot)
	return slot.ID()
}

// TestRecordMakeup 记录补课成功，预约即生效（StatusBooked），不需要审批。
func TestRecordMakeup(t *testing.T) {
	f := newFixture(t)
	cmd := f.validCmd()

	got, err := f.handler.RecordMakeup(context.Background(), cmd)
	if err != nil {
		t.Fatalf("RecordMakeup: %v", err)
	}
	if got.StudentID() != cmd.StudentID || got.CourseID() != cmd.CourseID {
		t.Errorf("makeup = (%d, %s), want (%d, %s)",
			got.StudentID(), got.CourseID(), cmd.StudentID, cmd.CourseID)
	}
	if got.Status() != makeup.StatusBooked {
		t.Errorf("status = %v, want %v", got.Status(), makeup.StatusBooked)
	}
	if got.MakeupHours() != 2 {
		t.Errorf("makeupHours = %d, want 2", got.MakeupHours())
	}
	if records := f.data.Makeups(); len(records) != 1 {
		t.Errorf("补课记录数 = %d, want 1", len(records))
	}
}

// TestRecordMakeupCountsOtherMakeups 教室装不装得下，要把已经约在同一节课上的其他
// 补课学员一起算：20 人的课 + 30 座教室，最多再坐 9 个补课学员（20+9+1=30）。
func TestRecordMakeupCountsOtherMakeups(t *testing.T) {
	cases := []struct {
		name     string
		existing int
		wantErr  bool
	}{
		{"前面 9 个，这一位刚好坐下", 9, false},
		{"前面 10 个，加上这一位就超了", 10, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			cmd := f.validCmd()
			f.seedMakeups(t, tc.existing, cmd.TargetSlotID, cmd.TargetDate, makeup.StatusBooked)

			_, err := f.handler.RecordMakeup(context.Background(), cmd)
			if tc.wantErr {
				if !errors.Is(err, makeup.ErrMakeupConflict) {
					t.Errorf("err = %v, want %v", err, makeup.ErrMakeupConflict)
				}
				return
			}
			if err != nil {
				t.Fatalf("RecordMakeup: %v", err)
			}
		})
	}
}

// TestRecordMakeupIgnoresOtherOccurrences 只有「同一槽位 + 同一日期 + 还等着上课」
// 的那些才占座：已取消的、以及同一个周排在别的日期上的，都不算。
func TestRecordMakeupIgnoresOtherOccurrences(t *testing.T) {
	f := newFixture(t)
	cmd := f.validCmd()

	// 同一节课上 10 条已取消 + 同一槽位下一个日期 10 条已预约 —— 都不该算进来
	f.seedMakeups(t, 10, cmd.TargetSlotID, cmd.TargetDate, makeup.StatusCancelled)
	f.seedMakeups(t, 10, cmd.TargetSlotID, cmd.TargetDate.AddDate(0, 0, 7), makeup.StatusBooked)

	if _, err := f.handler.RecordMakeup(context.Background(), cmd); err != nil {
		t.Fatalf("这些都不该占座，却报错 %v", err)
	}
}

// TestRecordMakeupClassroomTooSmall 目标那节课的教室装不下「课程人数上限 + 蹭课的这位」-> 拒绝。
func TestRecordMakeupClassroomTooSmall(t *testing.T) {
	f := newFixture(t)
	cmd := f.validCmd()
	cmd.TargetSlotID = f.smallSlot

	got, err := f.handler.RecordMakeup(context.Background(), cmd)
	if !errors.Is(err, makeup.ErrMakeupConflict) {
		t.Errorf("err = %v, want %v", err, makeup.ErrMakeupConflict)
	}
	if records := f.data.Makeups(); len(records) != 0 {
		t.Errorf("被拒绝时不应写入，补课记录数 = %d", len(records))
	}
	if got != nil {
		t.Errorf("被拒绝时不应返回记录，得到 %+v", got)
	}
}

// TestRecordMakeupRejected 校验不过时不落库。
func TestRecordMakeupRejected(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(cmd *RecordMakeup)
	}{
		{"学员 ID 非法", func(cmd *RecordMakeup) { cmd.StudentID = 0 }},
		{"课程为空", func(cmd *RecordMakeup) { cmd.CourseID = "" }},
		{"目标槽位为空", func(cmd *RecordMakeup) { cmd.TargetSlotID = "" }},
		{"补课课时非正", func(cmd *RecordMakeup) { cmd.MakeupHours = 0 }},
		// 目标槽位对不上真排期：服务里的 GetSlots 直接报 not found，不当作「装不下」
		{"目标槽位不存在", func(cmd *RecordMakeup) { cmd.TargetSlotID = "no-such-slot" }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			cmd := f.validCmd()
			tc.mutate(&cmd)

			got, err := f.handler.RecordMakeup(context.Background(), cmd)
			if err == nil {
				t.Fatalf("期望被拒绝，实际成功: %+v", got)
			}
			if records := f.data.Makeups(); len(records) != 0 {
				t.Errorf("被拒绝时不应写入，补课记录数 = %d", len(records))
			}
		})
	}
}

// TestCompleteAttendanceAfterRecord 记录后可以直接核销出勤，不需要先审批。
func TestCompleteAttendanceAfterRecord(t *testing.T) {
	f := newFixture(t)

	got, err := f.handler.RecordMakeup(context.Background(), f.validCmd())
	if err != nil {
		t.Fatalf("RecordMakeup: %v", err)
	}

	if err := got.CompleteAttendance(time.Now()); err != nil {
		t.Fatalf("CompleteAttendance: %v", err)
	}
	if got.Status() != makeup.StatusCompleted {
		t.Errorf("status = %v, want %v", got.Status(), makeup.StatusCompleted)
	}
	// 已完成的不能再取消
	if err := got.Cancel(101, time.Now()); !errors.Is(err, makeup.ErrAlreadyFinalized) {
		t.Errorf("Cancel(completed) = %v, want %v", err, makeup.ErrAlreadyFinalized)
	}
}

// TestDeleteMakeup 删除不存在的记录返回 not found。
func TestDeleteMakeup(t *testing.T) {
	f := newFixture(t)

	got, err := f.handler.RecordMakeup(context.Background(), f.validCmd())
	if err != nil {
		t.Fatalf("RecordMakeup: %v", err)
	}

	if err := f.handler.MakeupCmd.Delete(context.Background(), got.ID()); err != nil {
		t.Errorf("Delete(existing) = %v, want nil", err)
	}
	if err := f.handler.MakeupCmd.Delete(context.Background(), got.ID()); !errors.Is(err, makeup.ErrMakeupNotFound) {
		t.Errorf("Delete(missing) = %v, want %v", err, makeup.ErrMakeupNotFound)
	}
}
