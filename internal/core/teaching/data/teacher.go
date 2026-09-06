package data

import (
	"context"

	"cqrs/internal/core/teaching/biz/teacher"
	"cqrs/internal/shared/types"
)

type teacherRepo struct {
	data *Data
}

// NewTeacherRepo creates a new TeacherRepo instance.
func NewTeacherRepo(d *Data) teacher.TeacherRepo {
	return &teacherRepo{data: d}
}

func (r *teacherRepo) FindByID(ctx context.Context, id string) (*teacher.Teacher, error) {
	// TODO: implement database query
	return nil, nil
}

func (r *teacherRepo) FindByName(ctx context.Context, name string) (*teacher.Teacher, error) {
	// TODO: implement database query
	return nil, nil
}

func (r *teacherRepo) ListTeachers(ctx context.Context, opts ...types.ListOption) ([]*teacher.Teacher, error) {
	// TODO: implement database query
	return nil, nil
}

func (r *teacherRepo) CreateTeacher(ctx context.Context, t *teacher.Teacher) (*teacher.Teacher, error) {
	// TODO: implement database insert
	return t, nil
}

func (r *teacherRepo) UpdateTeacher(ctx context.Context, t *teacher.Teacher) (*teacher.Teacher, error) {
	// TODO: implement database update
	return t, nil
}

func (r *teacherRepo) DeleteTeacher(ctx context.Context, id string) error {
	// TODO: implement database soft delete
	return nil
}
