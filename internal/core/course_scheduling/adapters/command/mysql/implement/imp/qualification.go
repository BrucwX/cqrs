package imp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/adapters/command/mysql/implement/help"
	"cqrs/internal/core/course_scheduling/adapters/command/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type QualificationImp struct {
	data *mysql.Data
}

var _ repo.QualificationCommand = (*QualificationImp)(nil)

func NewQualificationImp(d *mysql.Data) repo.QualificationCommand {
	return &QualificationImp{data: d}
}

// GrantQualification 授予授课资质：够不够格由调用方判完，这里只落库。
func (c *QualificationImp) GrantQualification(ctx context.Context, q *qualification.Qualification) error {
	if q == nil {
		return qualification.ErrQualificationRequired
	}
	po, err := model.QualificationToPO(q)
	if err != nil {
		return err
	}

	_, err = c.data.Conn(ctx).ExecContext(ctx, `
INSERT INTO qualification
  (id, teacher_id, course_type_id, certified_at, status, updated_at)
VALUES (?, ?, ?, ?, ?, ?)`,
		po.ID, po.TeacherID, po.CourseTypeID, po.CertifiedAt, po.Status, po.UpdatedAt,
	)
	return err
}

// Delete 删除授课资质；不存在时报 ErrQualificationNotFound。
func (c *QualificationImp) Delete(ctx context.Context, id int64) error {
	res, err := c.data.Conn(ctx).ExecContext(ctx, `DELETE FROM qualification WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("%w: %d", qualification.ErrQualificationNotFound, id)
	}
	return nil
}

// MustGet 取授课资质聚合；不存在时报 ErrQualificationNotFound。
func (c *QualificationImp) MustGet(ctx context.Context, id int64) (qualification.Qualification, error) {
	po, err := help.ScanQualification(c.data.Conn(ctx).QueryRowContext(ctx, `
SELECT `+help.QualificationColumns+`
  FROM qualification
 WHERE id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return qualification.Qualification{}, fmt.Errorf("%w: %d", qualification.ErrQualificationNotFound, id)
	}
	if err != nil {
		return qualification.Qualification{}, err
	}
	do, err := model.QualificationToDO(po)
	if err != nil {
		return qualification.Qualification{}, err
	}
	return *do, nil
}

// GetQualifications 取该讲师持有的全部资质。
func (c *QualificationImp) GetQualifications(ctx context.Context, teacherID int64) ([]qualification.Qualification, error) {
	rows, err := c.data.Conn(ctx).QueryContext(ctx, `
SELECT `+help.QualificationColumns+`
  FROM qualification
 WHERE teacher_id = ?
 ORDER BY id`, teacherID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	out := make([]qualification.Qualification, 0)
	for rows.Next() {
		po, err := help.ScanQualification(rows)
		if err != nil {
			return nil, err
		}
		do, err := model.QualificationToDO(po)
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
