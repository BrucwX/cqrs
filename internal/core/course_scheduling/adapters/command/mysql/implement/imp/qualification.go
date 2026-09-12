package imp

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type QualificationImp struct {
	data *mysql.Data
}

var _ repo.QualificationCommand = (*QualificationImp)(nil)

func NewQualificationImp(d *mysql.Data) repo.QualificationCommand {
	return &QualificationImp{data: d}
}

func (c *QualificationImp) GrantQualification(ctx context.Context, q *qualification.Qualification) error {
	panic("implement me")
}

func (c *QualificationImp) Delete(id int64) error {
	panic("implement me")
}

func (c *QualificationImp) GetQualifications(teacherID int64) ([]qualification.Qualification, error) {
	panic("implement me")
}
