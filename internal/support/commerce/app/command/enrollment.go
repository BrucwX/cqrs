package command

import "cqrs/internal/support/commerce/domain/repo"

// EnrollmentUsecase is the enrollment usecase.
type EnrollmentUsecase struct {
	Repo repo.EnrollmentRepo
}

// NewEnrollmentUsecase creates a new EnrollmentUsecase.
func NewEnrollmentUsecase(r repo.EnrollmentRepo) *EnrollmentUsecase {
	return &EnrollmentUsecase{Repo: r}
}
