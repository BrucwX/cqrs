package course

import (
	"context"

	"cqrs/internal/shared/types"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// CourseRepo defines the course repository interface.
type CourseRepo interface {
	FindByID(ctx context.Context, id string) (*Course, error)
	ListCourses(ctx context.Context, opts ...types.ListOption) ([]*Course, error)
	CreateCourse(ctx context.Context, c *Course) (*Course, error)
	UpdateCourse(ctx context.Context, c *Course, mask *fieldmaskpb.FieldMask) (*Course, error)
	DeleteCourse(ctx context.Context, id string) error
}
