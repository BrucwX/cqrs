package qualification

import (
	"context"
	"errors"
	"testing"
	"time"

	"cqrs/internal/core/course_scheduling/adapters/memory"
	memorycmd "cqrs/internal/core/course_scheduling/adapters/memory/command"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

func newHandler(t *testing.T) (*Handler, *memory.Data) {
	t.Helper()

	d, cleanup, err := memory.NewData(nil)
	if err != nil {
		t.Fatalf("new data: %v", err)
	}
	t.Cleanup(cleanup)

	return NewHandler(memorycmd.NewQualificationCommand(d)), d
}

func newQualification(t *testing.T) *qualification.Qualification {
	t.Helper()

	now := time.Now()
	q, err := qualification.NewQualification(1, "ct-demo", now, now.AddDate(1, 0, 0))
	if err != nil {
		t.Fatalf("new qualification: %v", err)
	}
	return q
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
