package imp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/adapters/command/mysql/implement/help"
	"cqrs/internal/core/course_scheduling/adapters/command/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type ClassroomImp struct {
	data *mysql.Data
}

var _ repo.ClassroomCommand = (*ClassroomImp)(nil)

func NewClassroomImp(d *mysql.Data) repo.ClassroomCommand {
	return &ClassroomImp{data: d}
}

// Create 新增教室。
func (c *ClassroomImp) Create(ctx context.Context, cl *classroom.Classroom) error {
	if cl == nil {
		return repo.ErrClassroomRequired
	}
	po := model.ClassroomToPO(cl)

	_, err := c.data.Conn(ctx).ExecContext(ctx, `
INSERT INTO classroom (id, building, floor, room, capacity, allocated, status)
VALUES (?, ?, ?, ?, ?, ?, ?)`,
		po.ID, po.Building, po.Floor, po.Room, po.Capacity, po.Allocated, po.Status,
	)
	return err
}

// Update 按 ID 取出教室交给 updateFn 改，改完写回。
//
// 教室不存在时由 MustGet 报 ErrClassroomNotFound。
func (c *ClassroomImp) Update(
	ctx context.Context,
	id string,
	updateFn func(ctx context.Context, cl *classroom.Classroom) (*classroom.Classroom, error),
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
		return repo.ErrClassroomRequired
	}

	npo := model.ClassroomToPO(updated)
	_, err = c.data.Conn(ctx).ExecContext(ctx, `
UPDATE classroom
   SET building = ?, floor = ?, room = ?, capacity = ?, allocated = ?, status = ?
 WHERE id = ?`,
		npo.Building, npo.Floor, npo.Room, npo.Capacity, npo.Allocated, npo.Status,
		id,
	)
	return err
}

// Delete 删除教室；不存在时报 ErrClassroomNotFound。
func (c *ClassroomImp) Delete(ctx context.Context, id string) error {
	res, err := c.data.Conn(ctx).ExecContext(ctx, `DELETE FROM classroom WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("%w: %s", repo.ErrClassroomNotFound, id)
	}
	return nil
}

// MustGet 取教室聚合本身；不存在时报 ErrClassroomNotFound。
func (c *ClassroomImp) MustGet(ctx context.Context, id string) (classroom.Classroom, error) {
	po, err := help.ScanClassroom(c.data.Conn(ctx).QueryRowContext(ctx, `
SELECT `+help.ClassroomColumns+`
  FROM classroom
 WHERE id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return classroom.Classroom{}, fmt.Errorf("%w: %s", repo.ErrClassroomNotFound, id)
	}
	if err != nil {
		return classroom.Classroom{}, err
	}
	cl, err := model.ClassroomToDO(po)
	if err != nil {
		return classroom.Classroom{}, err
	}
	return *cl, nil
}
