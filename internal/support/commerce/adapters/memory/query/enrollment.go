package query

import (
	"cqrs/internal/support/commerce/adapters/memory"
	"cqrs/internal/support/commerce/domain/aggregate/enrollment"
)

type enrollmentQuery struct {
	data *memory.Data
}

// NewEnrollmentQuery creates a new EnrollmentQuery instance.
func NewEnrollmentQuery(d *memory.Data) *enrollmentQuery {
	return &enrollmentQuery{data: d}
}

func (q *enrollmentQuery) GetByID(id string) (*enrollment.Enrollment, error) {
	// TODO: implement database query
	return nil, nil
}

func (q *enrollmentQuery) List() ([]*enrollment.Enrollment, error) {
	// TODO: implement database query
	return nil, nil
}

func (q *enrollmentQuery) ListByStudentID(studentID string) ([]*enrollment.Enrollment, error) {
	// TODO: implement database query
	return nil, nil
}

func (q *enrollmentQuery) ListByCourseID(courseID string) ([]*enrollment.Enrollment, error) {
	// TODO: implement database query
	return nil, nil
}
