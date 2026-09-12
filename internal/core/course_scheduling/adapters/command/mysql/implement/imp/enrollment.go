package imp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/adapters/command/mysql/implement/help"
	"cqrs/internal/core/course_scheduling/adapters/command/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type CourseEnrollmentImp struct {
	data *mysql.Data
}

var _ repo.CourseEnrollmentCommand = (*CourseEnrollmentImp)(nil)

func NewCourseEnrollmentImp(d *mysql.Data) repo.CourseEnrollmentCommand {
	return &CourseEnrollmentImp{data: d}
}

// Save 保存课程注册（新增或更新，不做检查）。
func (c *CourseEnrollmentImp) Save(ctx context.Context, e *enrollment.CourseEnrollment) error {
	if e == nil {
		return enrollment.ErrEnrollmentRequired
	}
	po, err := model.EnrollmentToPO(e)
	if err != nil {
		return err
	}

	_, err = c.data.Conn(ctx).ExecContext(ctx, `
INSERT INTO course_enrollment
  (id, student_id, course_id, status, enrolled_at, completed_at, dropped_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  student_id = VALUES(student_id), course_id = VALUES(course_id),
  status = VALUES(status), enrolled_at = VALUES(enrolled_at),
  completed_at = VALUES(completed_at), dropped_at = VALUES(dropped_at),
  updated_at = VALUES(updated_at)`,
		po.ID, po.StudentID, po.CourseID, po.Status,
		po.EnrolledAt, po.CompletedAt, po.DroppedAt, po.UpdatedAt,
	)
	return err
}

// Delete 删除课程注册；不存在时报 ErrEnrollmentNotFound。
func (c *CourseEnrollmentImp) Delete(ctx context.Context, id int64) error {
	res, err := c.data.Conn(ctx).ExecContext(ctx, `DELETE FROM course_enrollment WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("%w: %d", enrollment.ErrEnrollmentNotFound, id)
	}
	return nil
}

// MustGet 取课程注册聚合；不存在时报 ErrEnrollmentNotFound。
func (c *CourseEnrollmentImp) MustGet(ctx context.Context, id int64) (enrollment.CourseEnrollment, error) {
	po, err := help.ScanEnrollment(c.data.Conn(ctx).QueryRowContext(ctx, `
SELECT `+help.EnrollmentColumns+`
  FROM course_enrollment
 WHERE id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return enrollment.CourseEnrollment{}, fmt.Errorf("%w: %d", enrollment.ErrEnrollmentNotFound, id)
	}
	if err != nil {
		return enrollment.CourseEnrollment{}, err
	}
	do, err := model.EnrollmentToDO(po)
	if err != nil {
		return enrollment.CourseEnrollment{}, err
	}
	return *do, nil
}

// Enroll 学员选课：准入判定由调用方做完，这里只落库。
func (c *CourseEnrollmentImp) Enroll(ctx context.Context, e *enrollment.CourseEnrollment) error {
	if e == nil {
		return enrollment.ErrEnrollmentRequired
	}
	return c.Save(ctx, e)
}

// GetEnrollments 取该学员的全部报名记录。
func (c *CourseEnrollmentImp) GetEnrollments(ctx context.Context, studentID int64) ([]enrollment.CourseEnrollment, error) {
	rows, err := c.data.Conn(ctx).QueryContext(ctx, `
SELECT `+help.EnrollmentColumns+`
  FROM course_enrollment
 WHERE student_id = ?
 ORDER BY id`, studentID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	out := make([]enrollment.CourseEnrollment, 0)
	for rows.Next() {
		po, err := help.ScanEnrollment(rows)
		if err != nil {
			return nil, err
		}
		do, err := model.EnrollmentToDO(po)
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
