package enrollment

import (
	"context"

	"github.com/go-kratos/kratos-layout/internal/shared/types"
)

// EnrollmentRepo defines the enrollment repository interface.
type EnrollmentRepo interface {
	FindByID(ctx context.Context, id string) (*Enrollment, error)
	FindByStudentAndCourse(ctx context.Context, studentID, courseID string) (*Enrollment, error)
	ListEnrollments(ctx context.Context, opts ...types.ListOption) ([]*Enrollment, error)
	CreateEnrollment(ctx context.Context, e *Enrollment) (*Enrollment, error)
	UpdateEnrollment(ctx context.Context, e *Enrollment) (*Enrollment, error)
	DeleteEnrollment(ctx context.Context, id string) error
}
