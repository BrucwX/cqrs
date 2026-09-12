package student

import (
	"context"
	"errors"
	"testing"
	"time"

	commandmemory "cqrs/internal/core/course_scheduling/adapters/command/memory"
	memorycmd "cqrs/internal/core/course_scheduling/adapters/command/memory/implement"
	"cqrs/internal/core/course_scheduling/adapters/memorystore"
	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
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
	return NewHandler(memorycmd.NewStudentCommand(data), data), store
}

func newStudent(t *testing.T) *student.Student {
	t.Helper()

	contact, err := student.NewContactInfo("13800000000", "")
	if err != nil {
		t.Fatalf("new contact: %v", err)
	}
	// ID 固定，所以走 Reconstitute 而不是会自己生成 ID 的 NewStudent
	return student.Reconstitute(1, "陈晨", student.TypeExternal, contact, student.StatusActive, time.Now(), time.Now())
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
