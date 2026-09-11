package student

import (
	"context"
	"errors"
	"testing"

	"cqrs/internal/core/course_scheduling/adapters/memory"
	memorycmd "cqrs/internal/core/course_scheduling/adapters/memory/command"
	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

func newHandler(t *testing.T) (*Handler, *memory.Data) {
	t.Helper()

	d, cleanup, err := memory.NewData(nil)
	if err != nil {
		t.Fatalf("new data: %v", err)
	}
	t.Cleanup(cleanup)

	return NewHandler(memorycmd.NewStudentCommand(d)), d
}

func newStudent(t *testing.T) *student.Student {
	t.Helper()

	contact, err := student.NewContactInfo("13800000000", "")
	if err != nil {
		t.Fatalf("new contact: %v", err)
	}
	s, err := student.NewStudent("陈晨", student.TypeExternal, contact)
	if err != nil {
		t.Fatalf("new student: %v", err)
	}
	return s
}

// TestDeleteStudent 删除存在的学员成功，再删报 not found。
func TestDeleteStudent(t *testing.T) {
	h, d := newHandler(t)
	s := newStudent(t)
	d.SeedStudent(s)

	if err := h.DeleteStudent(context.Background(), s.ID()); err != nil {
		t.Fatalf("DeleteStudent: %v", err)
	}
	if got := d.Students(); len(got) != 0 {
		t.Errorf("学员数量 = %d, want 0", len(got))
	}

	if err := h.DeleteStudent(context.Background(), s.ID()); !errors.Is(err, repo.ErrStudentNotFound) {
		t.Errorf("重复删除 err = %v, want %v", err, repo.ErrStudentNotFound)
	}
}
