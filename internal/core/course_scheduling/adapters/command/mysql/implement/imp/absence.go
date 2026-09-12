package imp

import (
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

func (c *AbsenceRecordImp) Save(a *absence.AbsenceRecord) error {
	panic("implement me")
}

func (c *AbsenceRecordImp) Delete(id int64) error {
	panic("implement me")
}
