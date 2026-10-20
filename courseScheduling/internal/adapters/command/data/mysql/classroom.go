package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"courseScheduling/internal/adapters/command/data/mysql/help"
	"courseScheduling/internal/adapters/command/data/mysql/model"
	appcmd "courseScheduling/internal/app/command/appRepo"
	"courseScheduling/internal/domain/aggregate/classroom"
)

// Create 新增教室。
func (c *MysqlData) CreateClassroom(ctx context.Context, cl *classroom.Classroom) error {
	if cl == nil {
		return classroom.ErrClassroomRequired
	}
	po, err := model.ClassroomToPO(cl)
	if err != nil {
		return err
	}

	_, err = c.Conn(ctx).ExecContext(ctx, `
INSERT INTO classroom (id, building, floor, room, capacity, allocated, status)
VALUES (?, ?, ?, ?, ?, ?, ?)`,
		po.ID, po.Building, po.Floor, po.Room, po.Capacity, po.Allocated, po.Status,
	)
	return err
}

// UpdateClassroom 按 ID 取出教室交给 updateFn 改，改完写回（带乐观锁版本校验）。
//
// 读出来的 lock_version 直接用来做 UPDATE 的 WHERE：对不上 ⇒ 影响 0 行 ⇒
// 返回 appcmd.ErrVersionConflict；写成才 lock_version + 1，
// 所以「读到 → 改 → 写回」不会悄悄盖掉别人的改动。
//
// ctx 里带着真事务时（重试耗尽后的兜底那条路）读带 FOR UPDATE 锁行，这次必然写成。
func (c *MysqlData) UpdateClassroom(
	ctx context.Context,
	id string,
	updateFn func(ctx context.Context, cl *classroom.Classroom) (*classroom.Classroom, error),
) error {
	po, err := c.classroomPO(ctx, id, c.inTx(ctx))
	if err != nil {
		return err
	}
	return c.updateClassroomPO(ctx, po, updateFn)
}

// updateClassroomPO 拿一个已经读出来的 PO（**带着读时的 lock_version**）走
// 「转 DO → updateFn → 带版本写回」。
//
// 单独抽出来是为了让乐观重试那一支把每轮新读的 PO 直接递进来 ——
// 版本始终跟着这一轮的读走，不用往 ctx 里塞。
func (c *MysqlData) updateClassroomPO(
	ctx context.Context,
	po *model.Classroom,
	updateFn func(ctx context.Context, cl *classroom.Classroom) (*classroom.Classroom, error),
) error {
	cl, err := model.ClassroomToDO(po)
	if err != nil {
		return err
	}

	updated, err := updateFn(ctx, cl)
	if err != nil {
		return err
	}
	if updated == nil {
		return classroom.ErrClassroomRequired
	}

	npo, err := model.ClassroomToPO(updated)
	if err != nil {
		return err
	}

	res, err := c.Conn(ctx).ExecContext(ctx, `
UPDATE classroom
   SET building = ?, floor = ?, room = ?, capacity = ?, allocated = ?, status = ?,
       lock_version = lock_version + 1
 WHERE id = ? AND lock_version = ?`,
		npo.Building, npo.Floor, npo.Room, npo.Capacity, npo.Allocated, npo.Status,
		po.ID, po.LockVersion,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("%w: classroom %s (version %d)", appcmd.ErrVersionConflict, po.ID, po.LockVersion)
	}
	return nil
}

// Delete 删除教室；不存在时报 ErrClassroomNotFound。
func (c *MysqlData) DeleteClassroom(ctx context.Context, id string) error {
	res, err := c.Conn(ctx).ExecContext(ctx, `DELETE FROM classroom WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("%w: %s", classroom.ErrClassroomNotFound, id)
	}
	return nil
}

// MustGetClassroom 取教室聚合本身；不存在时报 ErrClassroomNotFound。
func (c *MysqlData) MustGetClassroom(ctx context.Context, id string) (classroom.Classroom, error) {
	po, err := c.classroomPO(ctx, id, false)
	if err != nil {
		return classroom.Classroom{}, err
	}
	cl, err := model.ClassroomToDO(po)
	if err != nil {
		return classroom.Classroom{}, err
	}
	return *cl, nil
}

// classroomPO 读教室的 PO —— 比 MustGetClassroom 多带 lock_version，
// 写回时要用它做版本校验，所以不能走那个只返回聚合的读。
//
// forUpdate 只在真事务里开（见 UpdateClassroom）。
func (c *MysqlData) classroomPO(ctx context.Context, id string, forUpdate bool) (*model.Classroom, error) {
	query := `
SELECT ` + help.ClassroomColumns + `
  FROM classroom
 WHERE id = ?`
	if forUpdate {
		query += ` FOR UPDATE`
	}

	po, err := help.ScanClassroom(c.Conn(ctx).QueryRowContext(ctx, query, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: %s", classroom.ErrClassroomNotFound, id)
	}
	if err != nil {
		return nil, err
	}
	return po, nil
}
