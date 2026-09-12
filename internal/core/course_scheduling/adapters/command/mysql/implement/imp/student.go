package imp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/adapters/command/mysql/implement/help"
	"cqrs/internal/core/course_scheduling/adapters/command/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type StudentImp struct {
	data *mysql.Data
}

var _ repo.StudentCommand = (*StudentImp)(nil)

func NewStudentImp(d *mysql.Data) repo.StudentCommand {
	return &StudentImp{data: d}
}

// Create 新增学员。
func (c *StudentImp) Create(ctx context.Context, s *student.Student) error {
	if s == nil {
		return student.ErrStudentRequired
	}
	po, err := model.StudentToPO(s)
	if err != nil {
		return err
	}

	_, err = c.data.Conn(ctx).ExecContext(ctx, `
INSERT INTO student
  (id, name, student_type, phone, email, status, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		po.ID, po.Name, po.StudentType, po.Phone, po.Email, po.Status, po.CreatedAt, po.UpdatedAt,
	)
	return err
}

// Update 按 ID 取出学员交给 updateFn 改，改完写回。
func (c *StudentImp) Update(
	ctx context.Context,
	id int64,
	updateFn func(ctx context.Context, s *student.Student) (*student.Student, error),
) error {
	st, err := c.MustGet(ctx, id)
	if err != nil {
		return err
	}

	updated, err := updateFn(ctx, &st)
	if err != nil {
		return err
	}
	if updated == nil {
		return student.ErrStudentRequired
	}

	po, err := model.StudentToPO(updated)
	if err != nil {
		return err
	}
	_, err = c.data.Conn(ctx).ExecContext(ctx, `
UPDATE student
   SET name = ?, student_type = ?, phone = ?, email = ?, status = ?, updated_at = ?
 WHERE id = ?`,
		po.Name, po.StudentType, po.Phone, po.Email, po.Status, po.UpdatedAt, id,
	)
	return err
}

// Delete 删除学员；不存在时报 ErrStudentNotFound。
func (c *StudentImp) Delete(ctx context.Context, id int64) error {
	res, err := c.data.Conn(ctx).ExecContext(ctx, `DELETE FROM student WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("%w: %d", student.ErrStudentNotFound, id)
	}
	return nil
}

// Get 取学员；不存在时返回 (nil, nil)。
func (c *StudentImp) Get(ctx context.Context, id int64) (*student.Student, error) {
	do, err := c.MustGet(ctx, id)
	if errors.Is(err, student.ErrStudentNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &do, nil
}

// MustGet 取学员聚合本身；不存在时报 ErrStudentNotFound。
func (c *StudentImp) MustGet(ctx context.Context, id int64) (student.Student, error) {
	po, err := help.ScanStudent(c.data.Conn(ctx).QueryRowContext(ctx, `
SELECT `+help.StudentColumns+`
  FROM student
 WHERE id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return student.Student{}, fmt.Errorf("%w: %d", student.ErrStudentNotFound, id)
	}
	if err != nil {
		return student.Student{}, err
	}
	do, err := model.StudentToDO(po)
	if err != nil {
		return student.Student{}, err
	}
	return *do, nil
}
