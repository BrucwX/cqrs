package classroom

import (
	"context"
	"errors"
	"testing"

	"cqrs/internal/core/course_scheduling/adapters/memory"
	memorycmd "cqrs/internal/core/course_scheduling/adapters/memory/command"
	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

func newHandler(t *testing.T) (*Handler, *memory.Data) {
	t.Helper()

	d, cleanup, err := memory.NewData(nil)
	if err != nil {
		t.Fatalf("new data: %v", err)
	}
	t.Cleanup(cleanup)

	return NewHandler(memorycmd.NewClassroomCommand(d)), d
}

func newTestClassroom(t *testing.T, id string) *classroom.Classroom {
	t.Helper()

	loc, err := classroom.NewLocation("A", 1, id)
	if err != nil {
		t.Fatalf("new location: %v", err)
	}
	// ID 固定，所以走 Reconstitute 而不是会自己生成 ID 的 NewClassroom
	return classroom.Reconstitute(id, loc, 30, 0, classroom.StatusAvailable)
}

// TestDeleteClassroom 删除存在的教室成功，再删报 not found。
func TestDeleteClassroom(t *testing.T) {
	h, d := newHandler(t)
	room := newTestClassroom(t, "R101")
	d.SeedClassroom(room)

	if err := h.DeleteClassroom(context.Background(), room.ID()); err != nil {
		t.Fatalf("DeleteClassroom: %v", err)
	}
	if got := d.Classrooms(); len(got) != 0 {
		t.Errorf("教室数量 = %d, want 0", len(got))
	}

	if err := h.DeleteClassroom(context.Background(), room.ID()); !errors.Is(err, repo.ErrClassroomNotFound) {
		t.Errorf("重复删除 err = %v, want %v", err, repo.ErrClassroomNotFound)
	}
}
