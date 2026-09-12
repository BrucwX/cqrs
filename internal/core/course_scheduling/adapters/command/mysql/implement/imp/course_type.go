package imp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/adapters/command/mysql/implement/help"
	"cqrs/internal/core/course_scheduling/adapters/command/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type CourseTypeImp struct {
	data *mysql.Data
}

var _ repo.CourseTypeCommand = (*CourseTypeImp)(nil)

func NewCourseTypeImp(d *mysql.Data) repo.CourseTypeCommand {
	return &CourseTypeImp{data: d}
}

// MustGet 取课程类型聚合；不存在时报 ErrCourseTypeNotFound。
func (c *CourseTypeImp) MustGet(ctx context.Context, id string) (courseType.CourseType, error) {
	po, err := help.ScanCourseType(c.data.Conn(ctx).QueryRowContext(ctx, `
SELECT `+help.CourseTypeColumns+`
  FROM course_type
 WHERE id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return courseType.CourseType{}, fmt.Errorf("%w: %s", courseType.ErrCourseTypeNotFound, id)
	}
	if err != nil {
		return courseType.CourseType{}, err
	}
	do, err := model.CourseTypeToDO(po)
	if err != nil {
		return courseType.CourseType{}, err
	}
	return *do, nil
}

// GetCourseType 取某门课程归属的课程类型。
func (c *CourseTypeImp) GetCourseType(ctx context.Context, courseID string) (courseType.CourseType, error) {
	var courseTypeID string
	err := c.data.Conn(ctx).QueryRowContext(ctx, `
SELECT course_type_id
  FROM course
 WHERE id = ?`, courseID).Scan(&courseTypeID)
	if errors.Is(err, sql.ErrNoRows) {
		return courseType.CourseType{}, fmt.Errorf("%w: course %s", courseType.ErrCourseTypeNotFound, courseID)
	}
	if err != nil {
		return courseType.CourseType{}, err
	}
	return c.MustGet(ctx, courseTypeID)
}
