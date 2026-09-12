package imp

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type StudentImp struct {
	data *mysql.Data
}

var _ repo.StudentCommand = (*StudentImp)(nil)

func NewStudentImp(d *mysql.Data) repo.StudentCommand {
	return &StudentImp{data: d}
}

func (c *StudentImp) Create(ctx context.Context, s *student.Student) error {
	panic("implement me")
}

func (c *StudentImp) Update(
	ctx context.Context,
	id int64,
	updateFn func(ctx context.Context, s *student.Student) (*student.Student, error),
) error {
	panic("implement me")
}

func (c *StudentImp) Delete(ctx context.Context, id int64) error {
	panic("implement me")
}

func (c *StudentImp) Get(ctx context.Context, id int64) (*student.Student, error) {
	panic("implement me")
}
