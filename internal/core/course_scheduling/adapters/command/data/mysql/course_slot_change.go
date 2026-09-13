package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/command/data/mysql/help"
	"cqrs/internal/core/course_scheduling/adapters/command/data/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
)

// Save 保存课表变更（新增或更新，不做检查）。
func (c *MysqlData) SaveCourseSlotChange(ctx context.Context, csc *courseSlotChange.CourseSlotChange) error {
	if csc == nil {
		return courseSlotChange.ErrSlotChangeRequired
	}
	po, err := model.CourseSlotChangeToPO(csc)
	if err != nil {
		return err
	}

	_, err = c.Conn(ctx).ExecContext(ctx, `
INSERT INTO course_slot_change
  (id, course_id, applicant_id, change_type, original_slot_id, original_date,
   original_teacher_id, original_classroom_id, original_start_time, original_end_time,
   target_start_at, target_end_at, target_teacher_id, target_classroom_id, reason,
   created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  course_id = VALUES(course_id), applicant_id = VALUES(applicant_id),
  change_type = VALUES(change_type), original_slot_id = VALUES(original_slot_id),
  original_date = VALUES(original_date), original_teacher_id = VALUES(original_teacher_id),
  original_classroom_id = VALUES(original_classroom_id), original_start_time = VALUES(original_start_time),
  original_end_time = VALUES(original_end_time), target_start_at = VALUES(target_start_at),
  target_end_at = VALUES(target_end_at), target_teacher_id = VALUES(target_teacher_id),
  target_classroom_id = VALUES(target_classroom_id), reason = VALUES(reason),
  updated_at = VALUES(updated_at)`,
		po.ID, po.CourseID, po.ApplicantID, po.ChangeType,
		po.OriginalSlotID, po.OriginalDate, po.OriginalTeacherID, po.OriginalClassroomID,
		po.OriginalStartTime, po.OriginalEndTime,
		po.TargetStartAt, po.TargetEndAt, po.TargetTeacherID, po.TargetClassroomID,
		po.Reason, po.CreatedAt, po.UpdatedAt,
	)
	return err
}

// Delete 删除课表变更；不存在时报 ErrSlotChangeNotFound。
func (c *MysqlData) DeleteCourseSlotChange(ctx context.Context, id int64) error {
	res, err := c.Conn(ctx).ExecContext(ctx, `DELETE FROM course_slot_change WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("%w: %d", courseSlotChange.ErrSlotChangeNotFound, id)
	}
	return nil
}

// MustGet 取课表变更聚合；不存在时报 ErrSlotChangeNotFound。
func (c *MysqlData) MustGetCourseSlotChange(ctx context.Context, id int64) (courseSlotChange.CourseSlotChange, error) {
	po, err := help.ScanCourseSlotChange(c.Conn(ctx).QueryRowContext(ctx, `
SELECT `+help.CourseSlotChangeColumns+`
  FROM course_slot_change
 WHERE id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return courseSlotChange.CourseSlotChange{}, fmt.Errorf("%w: %d", courseSlotChange.ErrSlotChangeNotFound, id)
	}
	if err != nil {
		return courseSlotChange.CourseSlotChange{}, err
	}
	do, err := model.CourseSlotChangeToDO(po)
	if err != nil {
		return courseSlotChange.CourseSlotChange{}, err
	}
	return *do, nil
}

// Change 登记一次临时换课：冲突判定由调用方做完，这里只落库。
func (c *MysqlData) Change(ctx context.Context, csc *courseSlotChange.CourseSlotChange) error {
	if csc == nil {
		return courseSlotChange.ErrSlotChangeRequired
	}
	return c.SaveCourseSlotChange(ctx, csc)
}

// GetOtherSlotChanges 取除该变更单以外的全部换课记录。
func (c *MysqlData) GetOtherSlotChanges(ctx context.Context, id int64) ([]*courseSlotChange.CourseSlotChange, error) {
	rows, err := c.Conn(ctx).QueryContext(ctx, `
SELECT `+help.CourseSlotChangeColumns+`
  FROM course_slot_change
 WHERE id <> ?
 ORDER BY id`, id)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	out := make([]*courseSlotChange.CourseSlotChange, 0)
	for rows.Next() {
		po, err := help.ScanCourseSlotChange(rows)
		if err != nil {
			return nil, err
		}
		do, err := model.CourseSlotChangeToDO(po)
		if err != nil {
			return nil, err
		}
		out = append(out, do)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
