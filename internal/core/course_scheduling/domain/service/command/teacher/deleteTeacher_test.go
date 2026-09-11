package teacher

import (
	"context"
	"errors"
	"testing"
	"time"

	"cqrs/internal/core/course_scheduling/adapters/memory"
	memorycmd "cqrs/internal/core/course_scheduling/adapters/memory/command"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

func newHandler(t *testing.T) (*Handler, *memory.Data) {
	t.Helper()

	d, cleanup, err := memory.NewData(nil)
	if err != nil {
		t.Fatalf("new data: %v", err)
	}
	t.Cleanup(cleanup)

	return NewHandler(memorycmd.NewTeacherCommand(d)), d
}

func newTeacher(t *testing.T) *teacher.Teacher {
	t.Helper()

	contact, err := teacher.NewContactInfo("13800000000", "")
	if err != nil {
		t.Fatalf("new contact: %v", err)
	}
	// ID 固定，所以走 Reconstitute 而不是会自己生成 ID 的 NewTeacher
	return teacher.Reconstitute(1, 201, "李娜", "讲师", contact, teacher.StatusActive, time.Now(), time.Now())
}

// TestDeleteTeacher 删除存在的讲师成功，再删报 not found。
func TestDeleteTeacher(t *testing.T) {
	h, d := newHandler(t)
	tt := newTeacher(t)
	d.SeedTeacher(tt)

	if err := h.DeleteTeacher(context.Background(), tt.ID()); err != nil {
		t.Fatalf("DeleteTeacher: %v", err)
	}
	if got := d.Teachers(); len(got) != 0 {
		t.Errorf("讲师数量 = %d, want 0", len(got))
	}

	if err := h.DeleteTeacher(context.Background(), tt.ID()); !errors.Is(err, repo.ErrTeacherNotFound) {
		t.Errorf("重复删除 err = %v, want %v", err, repo.ErrTeacherNotFound)
	}
}
