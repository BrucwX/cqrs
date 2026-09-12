package imp

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type TeacherImp struct {
	data *mysql.Data
}

var _ repo.TeacherCommand = (*TeacherImp)(nil)

func NewTeacherImp(d *mysql.Data) repo.TeacherCommand {
	return &TeacherImp{data: d}
}

func (c *TeacherImp) Create(ctx context.Context, t *teacher.Teacher) error {
	panic("implement me")
}

func (c *TeacherImp) Update(
	ctx context.Context,
	id int64,
	updateFn func(ctx context.Context, t *teacher.Teacher) (*teacher.Teacher, error),
) error {
	panic("implement me")
}

func (c *TeacherImp) Delete(ctx context.Context, id int64) error {
	panic("implement me")
}

func (c *TeacherImp) Get(ctx context.Context, id int64) (*teacher.Teacher, error) {
	panic("implement me")
}

func (c *TeacherImp) MustGet(ctx context.Context, teacherID int64) (teacher.Teacher, error) {
	panic("implement me")
}
