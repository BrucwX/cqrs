package teacher

import (
	"context"
	"errors"
	"testing"
	"time"

	commandmemory "cqrs/internal/core/course_scheduling/adapters/command/memory"
	memorycmd "cqrs/internal/core/course_scheduling/adapters/command/memory/implement"
	"cqrs/internal/core/course_scheduling/adapters/memorystore"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
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
	return NewHandler(memorycmd.NewTeacherCommand(data)), store
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
