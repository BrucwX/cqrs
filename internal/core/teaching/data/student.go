package data

import (
	"context"

	"github.com/go-kratos/kratos-layout/internal/core/teaching/biz/student"
	"github.com/go-kratos/kratos-layout/internal/shared/types"
)

type studentRepo struct {
	data *Data
}

// NewStudentRepo creates a new StudentRepo instance.
func NewStudentRepo(d *Data) student.StudentRepo {
	return &studentRepo{data: d}
}

func (r *studentRepo) FindByID(ctx context.Context, id string) (*student.Student, error) {
	// TODO: implement database query
	return nil, nil
}

func (r *studentRepo) FindByPhone(ctx context.Context, phone string) (*student.Student, error) {
	// TODO: implement database query
	return nil, nil
}

func (r *studentRepo) ListStudents(ctx context.Context, opts ...types.ListOption) ([]*student.Student, error) {
	// TODO: implement database query
	return nil, nil
}

func (r *studentRepo) CreateStudent(ctx context.Context, s *student.Student) (*student.Student, error) {
	// TODO: implement database insert
	return s, nil
}

func (r *studentRepo) UpdateStudent(ctx context.Context, s *student.Student) (*student.Student, error) {
	// TODO: implement database update
	return s, nil
}

func (r *studentRepo) DeleteStudent(ctx context.Context, id string) error {
	// TODO: implement database soft delete
	return nil
}
