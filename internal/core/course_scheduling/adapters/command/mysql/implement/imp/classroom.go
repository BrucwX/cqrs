package imp

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type ClassroomImp struct {
	data *mysql.Data
}

var _ repo.ClassroomCommand = (*ClassroomImp)(nil)

func NewClassroomImp(d *mysql.Data) repo.ClassroomCommand {
	return &ClassroomImp{data: d}
}

func (c *ClassroomImp) Create(cl *classroom.Classroom) error {
	panic("implement me")
}

func (c *ClassroomImp) Update(
	ctx context.Context,
	id string,
	updateFn func(ctx context.Context, cl *classroom.Classroom) (*classroom.Classroom, error),
) error {
	panic("implement me")
}

func (c *ClassroomImp) Delete(id string) error {
	panic("implement me")
}

func (c *ClassroomImp) Get(id string) (*classroom.Classroom, error) {
	panic("implement me")
}
