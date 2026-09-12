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

func (c *CourseSlotChangeImp) Save(ctx context.Context, csc *courseSlotChange.CourseSlotChange) error {
	panic("implement me")
}

func (c *CourseSlotChangeImp) Delete(ctx context.Context, id int64) error {
	panic("implement me")
}

func (c *CourseSlotChangeImp) Change(ctx context.Context, csc *courseSlotChange.CourseSlotChange) error {
	panic("implement me")
}

func (c *CourseSlotChangeImp) GetOtherSlotChanges(ctx context.Context, id int64) ([]*courseSlotChange.CourseSlotChange, error) {
	panic("implement me")
}
