package classroom

import (
	"context"

	"github.com/go-kratos/kratos-layout/internal/shared/types"
)

// ClassroomRepo defines the classroom repository interface.
type ClassroomRepo interface {
	FindByID(ctx context.Context, id string) (*Classroom, error)
	FindByName(ctx context.Context, name string) (*Classroom, error)
	ListClassrooms(ctx context.Context, opts ...types.ListOption) ([]*Classroom, error)
	CreateClassroom(ctx context.Context, c *Classroom) (*Classroom, error)
	UpdateClassroom(ctx context.Context, c *Classroom) (*Classroom, error)
	DeleteClassroom(ctx context.Context, id string) error
}
