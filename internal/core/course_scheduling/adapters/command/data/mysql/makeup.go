package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"cqrs/internal/core/course_scheduling/adapters/command/data/mysql/help"
	"cqrs/internal/core/course_scheduling/adapters/command/data/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
)

// Save 保存补课预约（新增或更新）。
func (c *MysqlData) SaveMakeup(ctx context.Context, m *makeup.StudentMakeup) error {
	if m == nil {
		return makeup.ErrMakeupRequired
	}
	po, err := model.MakeupToPO(m)
	if err != nil {
		return err
	}

	_, err = c.Conn(ctx).ExecContext(ctx, `
INSERT INTO student_makeup
  (id, student_id, course_id, original_slot_id, original_date, target_slot_id, target_date,
   makeup_hours, status, completed_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  student_id = VALUES(student_id), course_id = VALUES(course_id),
  original_slot_id = VALUES(original_slot_id), original_date = VALUES(original_date),
  target_slot_id = VALUES(target_slot_id), target_date = VALUES(target_date),
  makeup_hours = VALUES(makeup_hours), status = VALUES(status),
  completed_at = VALUES(completed_at), updated_at = VALUES(updated_at)`,
		po.ID, po.StudentID, po.CourseID,
		po.OriginalSlotID, po.OriginalDate, po.TargetSlotID, po.TargetDate,
		po.MakeupHours, po.Status, po.CompletedAt, po.CreatedAt, po.UpdatedAt,
	)
	return err
}

// Delete 删除补课预约；不存在时报 ErrMakeupNotFound。
func (c *MysqlData) DeleteMakeup(ctx context.Context, id int64) error {
	res, err := c.Conn(ctx).ExecContext(ctx, `DELETE FROM student_makeup WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("%w: %d", makeup.ErrMakeupNotFound, id)
	}
	return nil
}

// MustGet 取补课预约聚合；不存在时报 ErrMakeupNotFound。
func (c *MysqlData) MustGetMakeup(ctx context.Context, id int64) (makeup.StudentMakeup, error) {
	po, err := help.ScanMakeup(c.Conn(ctx).QueryRowContext(ctx, `
SELECT `+help.MakeupColumns+`
  FROM student_makeup
 WHERE id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return makeup.StudentMakeup{}, fmt.Errorf("%w: %d", makeup.ErrMakeupNotFound, id)
	}
	if err != nil {
		return makeup.StudentMakeup{}, err
	}
	do, err := model.MakeupToDO(po)
	if err != nil {
		return makeup.StudentMakeup{}, err
	}
	return *do, nil
}

// GetMakeupsForTarget 取补到同一节课上的全部补课预约。
func (c *MysqlData) GetMakeupsForTarget(ctx context.Context, targetSlotID string, targetDate time.Time) ([]makeup.StudentMakeup, error) {
	rows, err := c.Conn(ctx).QueryContext(ctx, `
SELECT `+help.MakeupColumns+`
  FROM student_makeup
 WHERE target_slot_id = ? AND target_date = ?
 ORDER BY id`, targetSlotID, targetDate)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	out := make([]makeup.StudentMakeup, 0)
	for rows.Next() {
		po, err := help.ScanMakeup(rows)
		if err != nil {
			return nil, err
		}
		do, err := model.MakeupToDO(po)
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
