package command

import (
	"cqrs/internal/support/commerce/domain/repo/command"
	"cqrs/internal/support/commerce/domain/repo/query"
)

// EnrollmentUsecase is the enrollment usecase.
type EnrollmentUsecase struct {
	Query   query.EnrollmentQuery
	Command command.EnrollmentCommand
}

// NewEnrollmentUsecase creates a new EnrollmentUsecase.
func NewEnrollmentUsecase(q query.EnrollmentQuery, c command.EnrollmentCommand) *EnrollmentUsecase {
	return &EnrollmentUsecase{Query: q, Command: c}
}
