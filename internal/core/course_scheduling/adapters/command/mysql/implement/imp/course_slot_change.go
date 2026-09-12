package imp

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type CourseSlotChangeImp struct {
	data *mysql.Data
}

var _ repo.CourseSlotChangeCommand = (*CourseSlotChangeImp)(nil)

func NewCourseSlotChangeImp(d *mysql.Data) repo.CourseSlotChangeCommand {
	return &CourseSlotChangeImp{data: d}
}

func (c *CourseSlotChangeImp) Save(csc *courseSlotChange.CourseSlotChange) error {
	panic("implement me")
}

func (c *CourseSlotChangeImp) Delete(id int64) error {
	panic("implement me")
}

func (c *CourseSlotChangeImp) Change(ctx context.Context, csc *courseSlotChange.CourseSlotChange,
	checkConflictFn func(ctx context.Context, csc *courseSlotChange.CourseSlotChange) (bool, error)) error {
	panic("implement me")
}
