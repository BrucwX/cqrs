package imp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/adapters/command/mysql/implement/help"
	"cqrs/internal/core/course_scheduling/adapters/command/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type CourseImp struct {
	data *mysql.Data
}

var _ repo.CourseCommand = (*CourseImp)(nil)

func NewCourseImp(d *mysql.Data) repo.CourseCommand {
	return &CourseImp{data: d}
}

// Create 新增课程。
func (c *CourseImp) Create(ctx context.Context, crs *course.Course) error {
	if crs == nil {
		return repo.ErrCourseRequired
	}
	po, err := model.CourseToPO(crs)
	if err != nil {
		return err
	}

	_, err = c.data.Conn(ctx).ExecContext(ctx, `
INSERT INTO course
  (id, course_type_id, capacity_max, capacity_enrolled, enroll_start_at, enroll_end_at, drop_deadline,
   period_start_at, period_end_at, total_hours, completed_hours)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		po.ID, po.CourseTypeID, po.CapacityMax, po.CapacityEnrolled,
		po.EnrollStartAt, po.EnrollEndAt, po.DropDeadline,
		po.PeriodStartAt, po.PeriodEndAt, po.TotalHours, po.CompletedHours,
	)
	return err
}

// Update 按 ID 取出课程交给 updateFn 改，改完写回。
func (c *CourseImp) Update(
	ctx context.Context,
	id string,
	updateFn func(ctx context.Context, crs *course.Course) (*course.Course, error),
) error {
	cl, err := c.MustGet(ctx, id)
	if err != nil {
		return err
	}

	updated, err := updateFn(ctx, &cl)
	if err != nil {
		return err
	}
	if updated == nil {
		return repo.ErrCourseRequired
	}

	po, err := model.CourseToPO(updated)
	if err != nil {
		return err
	}
	_, err = c.data.Conn(ctx).ExecContext(ctx, `
UPDATE course
   SET course_type_id = ?, capacity_max = ?, capacity_enrolled = ?,
       enroll_start_at = ?, enroll_end_at = ?, drop_deadline = ?,
       period_start_at = ?, period_end_at = ?, total_hours = ?, completed_hours = ?
 WHERE id = ?`,
		po.CourseTypeID, po.CapacityMax, po.CapacityEnrolled,
		po.EnrollStartAt, po.EnrollEndAt, po.DropDeadline,
		po.PeriodStartAt, po.PeriodEndAt, po.TotalHours, po.CompletedHours,
		id,
	)
	return err
}

// Delete 删除课程；不存在时报 ErrCourseNotFound。
func (c *CourseImp) Delete(ctx context.Context, id string) error {
	res, err := c.data.Conn(ctx).ExecContext(ctx, `DELETE FROM course WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("%w: %s", repo.ErrCourseNotFound, id)
	}
	return nil
}

// Get 取课程；不存在时返回 (nil, nil)。
func (c *CourseImp) Get(ctx context.Context, id string) (*course.Course, error) {
	do, err := c.MustGet(ctx, id)
	if errors.Is(err, repo.ErrCourseNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &do, nil
}

// MustGet 取课程聚合本身；取不到报 ErrCourseNotFound。
func (c *CourseImp) MustGet(ctx context.Context, id string) (course.Course, error) {
	po, err := help.ScanCourse(c.data.Conn(ctx).QueryRowContext(ctx, `
SELECT `+help.CourseColumns+`
  FROM course
 WHERE id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return course.Course{}, fmt.Errorf("%w: %s", repo.ErrCourseNotFound, id)
	}
	if err != nil {
		return course.Course{}, err
	}
	do, err := model.CourseToDO(po)
	if err != nil {
		return course.Course{}, err
	}
	return *do, nil
}

// GetCourses 取某课程类型下的全部课程。
func (c *CourseImp) GetCourses(ctx context.Context, courseTypeID string) ([]course.Course, error) {
	rows, err := c.data.Conn(ctx).QueryContext(ctx, `
SELECT `+help.CourseColumns+`
  FROM course
 WHERE course_type_id = ?
 ORDER BY id`, courseTypeID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	out := make([]course.Course, 0)
	for rows.Next() {
		po, err := help.ScanCourse(rows)
		if err != nil {
			return nil, err
		}
		do, err := model.CourseToDO(po)
		if err != nil {
			return nil, err
		}
		out = append(out, *do)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
