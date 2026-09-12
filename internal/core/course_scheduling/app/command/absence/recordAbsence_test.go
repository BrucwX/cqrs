package absence

import (
	"context"
	"errors"
	"testing"
	"time"

	commandmemory "cqrs/internal/core/course_scheduling/adapters/command/memory"
	memorycmd "cqrs/internal/core/course_scheduling/adapters/command/memory/implement"
	"cqrs/internal/core/course_scheduling/adapters/memorystore"
	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
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
	return NewHandler(memorycmd.NewAbsenceCommand(data), data), store
}

func validCmd() RecordAbsence {
	return RecordAbsence{
		StudentID:    101,
		CourseID:     "C001",
		CourseSlotID: "550e8400-e29b-41d4-a716-446655440002",
		ScheduleDate: time.Date(2026, time.September, 16, 0, 0, 0, 0, time.Local),
		MissedHours:  2,
		AbsenceType:  absence.TypePersonalLeave,
		Reason:       "家中急事",
	}
}

// TestRecordAbsence 记录缺课成功，落库的就是一条事实。
func TestRecordAbsence(t *testing.T) {
	h, d := newHandler(t)
	cmd := validCmd()

	got, err := h.RecordAbsence(context.Background(), cmd)
	if err != nil {
		t.Fatalf("RecordAbsence: %v", err)
	}
	if got.StudentID() != cmd.StudentID || got.CourseID() != cmd.CourseID {
		t.Errorf("record = (%d, %s), want (%d, %s)",
			got.StudentID(), got.CourseID(), cmd.StudentID, cmd.CourseID)
	}
	if got.AbsenceType() != absence.TypePersonalLeave {
		t.Errorf("absenceType = %v, want %v", got.AbsenceType(), absence.TypePersonalLeave)
	}
	if got.MissedHours() != 2 {
		t.Errorf("missedHours = %d, want 2", got.MissedHours())
	}
	if records := d.Absences(); len(records) != 1 {
		t.Errorf("缺勤记录数 = %d, want 1", len(records))
	}
}

// TestRecordAbsenceRejected 校验不过时不落库。
func TestRecordAbsenceRejected(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(cmd *RecordAbsence)
		wantErr error
	}{
		{"学员 ID 非法", func(cmd *RecordAbsence) { cmd.StudentID = 0 }, nil},
		{"槽位 ID 非法", func(cmd *RecordAbsence) { cmd.CourseSlotID = "" }, nil},
		{"课程为空", func(cmd *RecordAbsence) { cmd.CourseID = "" }, nil},
		{"课时数非正", func(cmd *RecordAbsence) { cmd.MissedHours = 0 }, absence.ErrInvalidHours},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h, d := newHandler(t)
			cmd := validCmd()
			tc.mutate(&cmd)

			got, err := h.RecordAbsence(context.Background(), cmd)
			if err == nil {
				t.Fatalf("期望被拒绝，实际成功: %+v", got)
			}
			if tc.wantErr != nil && !errors.Is(err, tc.wantErr) {
				t.Errorf("err = %v, want %v", err, tc.wantErr)
			}
			if records := d.Absences(); len(records) != 0 {
				t.Errorf("被拒绝时不应写入，缺勤记录数 = %d", len(records))
			}
		})
	}
}

// TestDeleteAbsence 删除不存在的记录返回 not found。
func TestDeleteAbsence(t *testing.T) {
	h, _ := newHandler(t)

	got, err := h.RecordAbsence(context.Background(), validCmd())
	if err != nil {
		t.Fatalf("RecordAbsence: %v", err)
	}

	if err := h.AbsenceCmd.Delete(context.Background(), got.ID()); err != nil {
		t.Errorf("Delete(existing) = %v, want nil", err)
	}
	if err := h.AbsenceCmd.Delete(context.Background(), got.ID()); !errors.Is(err, repo.ErrAbsenceNotFound) {
		t.Errorf("Delete(missing) = %v, want %v", err, repo.ErrAbsenceNotFound)
	}
}
