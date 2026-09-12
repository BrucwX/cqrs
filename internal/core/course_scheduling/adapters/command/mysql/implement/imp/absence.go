package imp

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type AbsenceRecordImp struct {
	data *mysql.Data
}

var _ repo.AbsenceRecordCommand = (*AbsenceRecordImp)(nil)

func NewAbsenceRecordImp(d *mysql.Data) repo.AbsenceRecordCommand {
	return &AbsenceRecordImp{data: d}
}

func (c *AbsenceRecordImp) Save(ctx context.Context, a *absence.AbsenceRecord) error {
	panic("implement me")
}

func (c *AbsenceRecordImp) Delete(ctx context.Context, id int64) error {
	panic("implement me")
}

func (c *AbsenceRecordImp) GetAbsences(ctx context.Context, studentID int64) ([]absence.AbsenceRecord, error) {
	panic("implement me")
}
