package qualification

import (
	"context"
	"errors"
	"testing"
	"time"

	commandmemory "cqrs/internal/core/course_scheduling/adapters/command/memory"
	memorycmd "cqrs/internal/core/course_scheduling/adapters/command/memory/implement"
	"cqrs/internal/core/course_scheduling/adapters/memorystore"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
	"cqrs/internal/core/course_scheduling/domain/service/qualificationCheck"
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
	return NewHandler(
		memorycmd.NewQualificationCommand(data),
		qualificationCheck.NewService(
			memorycmd.NewCourseSlotCommand(data),
			memorycmd.NewTeacherCommand(data),
			memorycmd.NewQualificationCommand(data),
			memorycmd.NewCourseTypeCommand(data),
			memorycmd.NewCourseCommand(data),
			memorycmd.NewEnrollmentCommand(data),
			memorycmd.NewAbsenceCommand(data),
		),
		data,
	), store
}

func newQualification(t *testing.T) *qualification.Qualification {
	t.Helper()

	now := time.Now()
	// ID 固定，所以走 Reconstitute 而不是会自己生成 ID 的 NewQualification
	return qualification.Reconstitute(
		1, 1, "ct-demo", now, qualification.StatusActive, now,
	)
}

// TestDeleteQualification 删除存在的资质成功，再删报 not found。
func TestDeleteQualification(t *testing.T) {
	h, d := newHandler(t)
	q := newQualification(t)
	d.SeedQualification(q)

	if err := h.DeleteQualification(context.Background(), q.ID()); err != nil {
		t.Fatalf("DeleteQualification: %v", err)
	}
	if got := d.Qualifications(); len(got) != 0 {
		t.Errorf("资质数量 = %d, want 0", len(got))
	}

	if err := h.DeleteQualification(context.Background(), q.ID()); !errors.Is(err, repo.ErrQualificationNotFound) {
		t.Errorf("重复删除 err = %v, want %v", err, repo.ErrQualificationNotFound)
	}
}
