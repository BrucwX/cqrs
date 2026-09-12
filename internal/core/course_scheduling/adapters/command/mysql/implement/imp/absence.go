package imp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/adapters/command/mysql/implement/help"
	"cqrs/internal/core/course_scheduling/adapters/command/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type AbsenceRecordImp struct {
	data *mysql.Data
}

var _ repo.AbsenceRecordCommand = (*AbsenceRecordImp)(nil)

func NewAbsenceRecordImp(d *mysql.Data) repo.AbsenceRecordCommand {
	return &AbsenceRecordImp{data: d}
}

// Save 保存缺勤记录（新增或更新）。
func (c *AbsenceRecordImp) Save(ctx context.Context, a *absence.AbsenceRecord) error {
	if a == nil {
		return absence.ErrAbsenceRequired
	}
	po, err := model.AbsenceToPO(a)
	if err != nil {
		return err
	}

	_, err = c.data.Conn(ctx).ExecContext(ctx, `
INSERT INTO absence_record
  (id, student_id, course_id, course_slot_id, schedule_date, missed_hours, absence_type, reason, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  student_id = VALUES(student_id), course_id = VALUES(course_id),
  course_slot_id = VALUES(course_slot_id), schedule_date = VALUES(schedule_date),
  missed_hours = VALUES(missed_hours), absence_type = VALUES(absence_type),
  reason = VALUES(reason), updated_at = VALUES(updated_at)`,
		po.ID, po.StudentID, po.CourseID, po.CourseSlotID, po.ScheduleDate,
		po.MissedHours, po.AbsenceType, po.Reason, po.CreatedAt, po.UpdatedAt,
	)
	return err
}

// Delete 删除缺勤记录；不存在时报 ErrAbsenceNotFound。
func (c *AbsenceRecordImp) Delete(ctx context.Context, id int64) error {
	res, err := c.data.Conn(ctx).ExecContext(ctx, `DELETE FROM absence_record WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("%w: %d", absence.ErrAbsenceNotFound, id)
	}
	return nil
}

// MustGet 取缺勤记录聚合；不存在时报 ErrAbsenceNotFound。
func (c *AbsenceRecordImp) MustGet(ctx context.Context, id int64) (absence.AbsenceRecord, error) {
	po, err := help.ScanAbsence(c.data.Conn(ctx).QueryRowContext(ctx, `
SELECT `+help.AbsenceColumns+`
  FROM absence_record
 WHERE id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return absence.AbsenceRecord{}, fmt.Errorf("%w: %d", absence.ErrAbsenceNotFound, id)
	}
	if err != nil {
		return absence.AbsenceRecord{}, err
	}
	do, err := model.AbsenceToDO(po)
	if err != nil {
		return absence.AbsenceRecord{}, err
	}
	return *do, nil
}

// GetAbsences 取该学员的全部缺勤记录。
func (c *AbsenceRecordImp) GetAbsences(ctx context.Context, studentID int64) ([]absence.AbsenceRecord, error) {
	rows, err := c.data.Conn(ctx).QueryContext(ctx, `
SELECT `+help.AbsenceColumns+`
  FROM absence_record
 WHERE student_id = ?
 ORDER BY id`, studentID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	out := make([]absence.AbsenceRecord, 0)
	for rows.Next() {
		po, err := help.ScanAbsence(rows)
		if err != nil {
			return nil, err
		}
		do, err := model.AbsenceToDO(po)
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
