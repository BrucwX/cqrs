package classroom

import (
	"context"
	"errors"
	"testing"

	commandmemory "cqrs/internal/core/course_scheduling/adapters/command/memory"
	memorycmd "cqrs/internal/core/course_scheduling/adapters/command/memory/implement"
	"cqrs/internal/core/course_scheduling/adapters/memorystore"
	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
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
	return NewHandler(memorycmd.NewClassroomCommand(data)), store
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
