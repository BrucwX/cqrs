package imp

import (
	"time"

	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type StudentMakeupImp struct {
	data *mysql.Data
}

var _ repo.StudentMakeupCommand = (*StudentMakeupImp)(nil)

func NewStudentMakeupImp(d *mysql.Data) repo.StudentMakeupCommand {
	return &StudentMakeupImp{data: d}
}

func (c *StudentMakeupImp) Save(m *makeup.StudentMakeup) error {
	panic("implement me")
}

func (c *StudentMakeupImp) Delete(id int64) error {
	panic("implement me")
}

func (c *StudentMakeupImp) GetMakeupsForTarget(targetSlotID string, targetDate time.Time) ([]makeup.StudentMakeup, error) {
	panic("implement me")
}
