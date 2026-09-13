package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/command/data/mysql/help"
	"cqrs/internal/core/course_scheduling/adapters/command/data/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
)

// MustGet 取课程类型聚合；不存在时报 ErrCourseTypeNotFound。
func (c *MysqlData) MustGetCourseType(ctx context.Context, id string) (courseType.CourseType, error) {
	po, err := help.ScanCourseType(c.Conn(ctx).QueryRowContext(ctx, `
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
func (c *MysqlData) GetCourseType(ctx context.Context, courseID string) (courseType.CourseType, error) {
	var courseTypeID string
	err := c.Conn(ctx).QueryRowContext(ctx, `
SELECT course_type_id
  FROM course
 WHERE id = ?`, courseID).Scan(&courseTypeID)
	if errors.Is(err, sql.ErrNoRows) {
		return courseType.CourseType{}, fmt.Errorf("%w: course %s", courseType.ErrCourseTypeNotFound, courseID)
	}
	if err != nil {
		return courseType.CourseType{}, err
	}
	return c.MustGetCourseType(ctx, courseTypeID)
}
