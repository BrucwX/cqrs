package command

import (
	"cqrs/internal/support/commerce/adapters/memory"
	"cqrs/internal/support/commerce/domain/aggregate/enrollment"
)

type enrollmentCommand struct {
	data *memory.Data
}

// NewEnrollmentCommand creates a new EnrollmentCommand instance.
func NewEnrollmentCommand(d *memory.Data) *enrollmentCommand {
	return &enrollmentCommand{data: d}
}

func (r *enrollmentCommand) Save(e *enrollment.Enrollment) error {
	// TODO: implement database insert/update
	return nil
}

func (r *enrollmentCommand) Delete(id string) error {
	// TODO: implement database soft delete
	return nil
}
