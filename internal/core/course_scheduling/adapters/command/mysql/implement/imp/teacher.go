package imp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/adapters/command/mysql/implement/help"
	"cqrs/internal/core/course_scheduling/adapters/command/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type TeacherImp struct {
	data *mysql.Data
}

var _ repo.TeacherCommand = (*TeacherImp)(nil)

func NewTeacherImp(d *mysql.Data) repo.TeacherCommand {
	return &TeacherImp{data: d}
}

// Create 新增讲师。
func (c *TeacherImp) Create(ctx context.Context, t *teacher.Teacher) error {
	if t == nil {
		return teacher.ErrTeacherRequired
	}
	po, err := model.TeacherToPO(t)
	if err != nil {
		return err
	}

	_, err = c.data.Conn(ctx).ExecContext(ctx, `
INSERT INTO teacher
  (id, student_id, name, title, phone, email, status, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		po.ID, po.StudentID, po.Name, po.Title, po.Phone, po.Email, po.Status, po.CreatedAt, po.UpdatedAt,
	)
	return err
}

// Update 按 ID 取出讲师交给 updateFn 改，改完写回。
func (c *TeacherImp) Update(
	ctx context.Context,
	id int64,
	updateFn func(ctx context.Context, t *teacher.Teacher) (*teacher.Teacher, error),
) error {
	t, err := c.MustGet(ctx, id)
	if err != nil {
		return err
	}

	updated, err := updateFn(ctx, &t)
	if err != nil {
		return err
	}
	if updated == nil {
		return teacher.ErrTeacherRequired
	}

	po, err := model.TeacherToPO(updated)
	if err != nil {
		return err
	}
	_, err = c.data.Conn(ctx).ExecContext(ctx, `
UPDATE teacher
   SET student_id = ?, name = ?, title = ?, phone = ?, email = ?, status = ?, updated_at = ?
 WHERE id = ?`,
		po.StudentID, po.Name, po.Title, po.Phone, po.Email, po.Status, po.UpdatedAt, id,
	)
	return err
}

// Delete 删除讲师；不存在时报 ErrTeacherNotFound。
func (c *TeacherImp) Delete(ctx context.Context, id int64) error {
	res, err := c.data.Conn(ctx).ExecContext(ctx, `DELETE FROM teacher WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("%w: %d", teacher.ErrTeacherNotFound, id)
	}
	return nil
}

// Get 取讲师；不存在时返回 (nil, nil)。
func (c *TeacherImp) Get(ctx context.Context, id int64) (*teacher.Teacher, error) {
	do, err := c.MustGet(ctx, id)
	if errors.Is(err, teacher.ErrTeacherNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &do, nil
}

// MustGet 取讲师聚合本身；不存在时报 ErrTeacherNotFound。
func (c *TeacherImp) MustGet(ctx context.Context, id int64) (teacher.Teacher, error) {
	po, err := help.ScanTeacher(c.data.Conn(ctx).QueryRowContext(ctx, `
SELECT `+help.TeacherColumns+`
  FROM teacher
 WHERE id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return teacher.Teacher{}, fmt.Errorf("%w: %d", teacher.ErrTeacherNotFound, id)
	}
	if err != nil {
		return teacher.Teacher{}, err
	}
	do, err := model.TeacherToDO(po)
	if err != nil {
		return teacher.Teacher{}, err
	}
	return *do, nil
}
