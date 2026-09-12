package makeup

import (
	"context"
	"errors"
	"testing"
	"time"

	commandmemory "cqrs/internal/core/course_scheduling/adapters/command/memory"
	memorycmd "cqrs/internal/core/course_scheduling/adapters/command/memory/implement"
	"cqrs/internal/core/course_scheduling/adapters/memorystore"
	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// newHandler 装配一个跑在干净内存存储上的命令处理器。
//
// 返回完整 store 给测试塞数据/做断言，仓库拿到的则是收窄后的写侧面。
func newHandler(t *testing.T) (*Handler, *memorystore.Data) {
	t.Helper()

	store, cleanup, err := memorystore.NewData(nil)
	if err != nil {
		t.Fatalf("new data: %v", err)
	}
	t.Cleanup(cleanup)

	data := commandmemory.NewData(store)
	return NewHandler(memorycmd.NewMakeupCommand(data)), store
}

func validCmd() RecordMakeup {
	return RecordMakeup{
		StudentID:      101,
		CourseID:       "C001",
		OriginalSlotID: 1,
		OriginalDate:   time.Date(2026, time.September, 14, 0, 0, 0, 0, time.Local),
		TargetSlotID:   2,
		TargetDate:     time.Date(2026, time.September, 16, 0, 0, 0, 0, time.Local),
		MakeupHours:    2,
	}
}

// TestRecordMakeup 记录补课成功，预约即生效（StatusBooked），不需要审批。
func TestRecordMakeup(t *testing.T) {
	h, d := newHandler(t)
	cmd := validCmd()

	got, err := h.RecordMakeup(context.Background(), cmd)
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
	if records := d.Makeups(); len(records) != 1 {
		t.Errorf("补课记录数 = %d, want 1", len(records))
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
		{"目标槽位为空", func(cmd *RecordMakeup) { cmd.TargetSlotID = 0 }},
		{"补课课时非正", func(cmd *RecordMakeup) { cmd.MakeupHours = 0 }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h, d := newHandler(t)
			cmd := validCmd()
			tc.mutate(&cmd)

			got, err := h.RecordMakeup(context.Background(), cmd)
			if err == nil {
				t.Fatalf("期望被拒绝，实际成功: %+v", got)
			}
			if records := d.Makeups(); len(records) != 0 {
				t.Errorf("被拒绝时不应写入，补课记录数 = %d", len(records))
			}
		})
	}
}

// TestCompleteAttendanceAfterRecord 记录后可以直接核销出勤，不需要先审批。
func TestCompleteAttendanceAfterRecord(t *testing.T) {
	h, _ := newHandler(t)

	got, err := h.RecordMakeup(context.Background(), validCmd())
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
	h, _ := newHandler(t)

	got, err := h.RecordMakeup(context.Background(), validCmd())
	if err != nil {
		t.Fatalf("RecordMakeup: %v", err)
	}

	if err := h.MakeupCmd.Delete(got.ID()); err != nil {
		t.Errorf("Delete(existing) = %v, want nil", err)
	}
	if err := h.MakeupCmd.Delete(got.ID()); !errors.Is(err, repo.ErrMakeupNotFound) {
		t.Errorf("Delete(missing) = %v, want %v", err, repo.ErrMakeupNotFound)
	}
}
