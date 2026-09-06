package data

import (
	"context"

	"github.com/go-kratos/kratos-layout/internal/core/venue/biz/classroom"
	"github.com/go-kratos/kratos-layout/internal/shared/types"
)

type classroomRepo struct {
	data *Data
}

// NewClassroomRepo creates a new ClassroomRepo instance.
func NewClassroomRepo(d *Data) classroom.ClassroomRepo {
	return &classroomRepo{data: d}
}

func (r *classroomRepo) FindByID(ctx context.Context, id string) (*classroom.Classroom, error) {
	// TODO: implement database query
	return nil, nil
}

func (r *classroomRepo) FindByName(ctx context.Context, name string) (*classroom.Classroom, error) {
	// TODO: implement database query
	return nil, nil
}

func (r *classroomRepo) ListClassrooms(ctx context.Context, opts ...types.ListOption) ([]*classroom.Classroom, error) {
	// TODO: implement database query
	return nil, nil
}

func (r *classroomRepo) CreateClassroom(ctx context.Context, c *classroom.Classroom) (*classroom.Classroom, error) {
	// TODO: implement database insert
	return c, nil
}

func (r *classroomRepo) UpdateClassroom(ctx context.Context, c *classroom.Classroom) (*classroom.Classroom, error) {
	// TODO: implement database update
	return c, nil
}

func (r *classroomRepo) DeleteClassroom(ctx context.Context, id string) error {
	// TODO: implement database soft delete
	return nil
}
