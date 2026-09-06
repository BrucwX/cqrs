package student

import (
	"context"

	"github.com/go-kratos/kratos-layout/internal/shared/types"
)

// StudentRepo defines the student repository interface.
type StudentRepo interface {
	FindByID(ctx context.Context, id string) (*Student, error)
	FindByPhone(ctx context.Context, phone string) (*Student, error)
	ListStudents(ctx context.Context, opts ...types.ListOption) ([]*Student, error)
	CreateStudent(ctx context.Context, s *Student) (*Student, error)
	UpdateStudent(ctx context.Context, s *Student) (*Student, error)
	DeleteStudent(ctx context.Context, id string) error
}
