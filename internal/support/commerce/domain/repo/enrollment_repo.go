package repo

import (
	"context"

	"cqrs/internal/support/commerce/domain/aggregate/enrollment"
)

// EnrollmentRepo defines the enrollment repository interface.
type EnrollmentRepo interface {
	FindByID(ctx context.Context, id string) (*enrollment.Enrollment, error)
	FindByStudentAndCourse(ctx context.Context, studentID, courseID string) (*enrollment.Enrollment, error)
	CreateEnrollment(ctx context.Context, e *enrollment.Enrollment) (*enrollment.Enrollment, error)
	UpdateEnrollment(ctx context.Context, e *enrollment.Enrollment) (*enrollment.Enrollment, error)
	DeleteEnrollment(ctx context.Context, id string) error
}
