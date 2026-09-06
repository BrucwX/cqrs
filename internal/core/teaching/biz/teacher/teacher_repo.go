package teacher

import (
	"context"

	"github.com/go-kratos/kratos-layout/internal/shared/types"
)

// TeacherRepo defines the teacher repository interface.
type TeacherRepo interface {
	FindByID(ctx context.Context, id string) (*Teacher, error)
	FindByName(ctx context.Context, name string) (*Teacher, error)
	ListTeachers(ctx context.Context, opts ...types.ListOption) ([]*Teacher, error)
	CreateTeacher(ctx context.Context, t *Teacher) (*Teacher, error)
	UpdateTeacher(ctx context.Context, t *Teacher) (*Teacher, error)
	DeleteTeacher(ctx context.Context, id string) error
}
