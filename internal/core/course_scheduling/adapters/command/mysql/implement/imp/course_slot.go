package imp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/adapters/command/mysql/implement/help"
	"cqrs/internal/core/course_scheduling/adapters/command/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type CourseSlotImp struct {
	data *mysql.Data
}

var _ repo.CourseSlotCommand = (*CourseSlotImp)(nil)

func NewCourseSlotImp(d *mysql.Data) repo.CourseSlotCommand {
	return &CourseSlotImp{data: d}
}

// Save 保存课表槽位（新增或更新）。
func (c *CourseSlotImp) Save(ctx context.Context, cs *courseSlot.CourseSlot) error {
	if cs == nil {
		return courseSlot.ErrCourseSlotRequired
	}
	po, err := model.CourseSlotToPO(cs)
	if err != nil {
		return err
	}

	_, err = c.data.Conn(ctx).ExecContext(ctx, `
INSERT INTO course_slot
  (id, course_id, weekday, start_time, end_time, teacher_id, classroom_id, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
ON DUPLICATE KEY UPDATE
  course_id = VALUES(course_id), weekday = VALUES(weekday),
  start_time = VALUES(start_time), end_time = VALUES(end_time),
  teacher_id = VALUES(teacher_id), classroom_id = VALUES(classroom_id),
  updated_at = VALUES(updated_at)`,
		po.ID, po.CourseID, po.Weekday, po.StartTime, po.EndTime,
		po.TeacherID, po.ClassroomID, po.CreatedAt, po.UpdatedAt,
	)
	return err
}

// Delete 删除课表槽位；不存在时报 ErrCourseSlotNotFound。
func (c *CourseSlotImp) Delete(ctx context.Context, id string) error {
	res, err := c.data.Conn(ctx).ExecContext(ctx, `DELETE FROM course_slot WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("%w: %s", courseSlot.ErrCourseSlotNotFound, id)
	}
	return nil
}

// MustGet 取课表槽位聚合；不存在时报 ErrCourseSlotNotFound。
func (c *CourseSlotImp) MustGet(ctx context.Context, id string) (courseSlot.CourseSlot, error) {
	po, err := help.ScanCourseSlot(c.data.Conn(ctx).QueryRowContext(ctx, `
SELECT `+help.CourseSlotColumns+`
  FROM course_slot
 WHERE id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return courseSlot.CourseSlot{}, fmt.Errorf("%w: %s", courseSlot.ErrCourseSlotNotFound, id)
	}
	if err != nil {
		return courseSlot.CourseSlot{}, err
	}
	do, err := model.CourseSlotToDO(po)
	if err != nil {
		return courseSlot.CourseSlot{}, err
	}
	return *do, nil
}

// AssignTeacher 给指定课表槽位们配置老师。
func (c *CourseSlotImp) AssignTeacher(ctx context.Context, slotIDs []string, teacherID int64) error {
	for _, id := range slotIDs {
		slot, err := c.MustGet(ctx, id)
		if err != nil {
			return err
		}
		if err := slot.ChangeTeacher(teacherID); err != nil {
			return err
		}
		if err := c.Save(ctx, &slot); err != nil {
			return err
		}
	}
	return nil
}

// AssignCourse 给指定课表槽位们配置课程。
func (c *CourseSlotImp) AssignCourse(ctx context.Context, slotIDs []string, courseID string) error {
	for _, id := range slotIDs {
		slot, err := c.MustGet(ctx, id)
		if err != nil {
			return err
		}
		if err := slot.ChangeCourse(courseID); err != nil {
			return err
		}
		if err := c.Save(ctx, &slot); err != nil {
			return err
		}
	}
	return nil
}

// AssignClassroom 给指定课表槽位们配置教室。
func (c *CourseSlotImp) AssignClassroom(ctx context.Context, slotIDs []string, classroomID string) error {
	for _, id := range slotIDs {
		slot, err := c.MustGet(ctx, id)
		if err != nil {
			return err
		}
		if err := slot.ChangeClassroom(classroomID); err != nil {
			return err
		}
		if err := c.Save(ctx, &slot); err != nil {
			return err
		}
	}
	return nil
}

// GetSlots 按 ID 取槽位，顺序与入参一致；少一个就报 not found。
func (c *CourseSlotImp) GetSlots(ctx context.Context, slotIDs []string) (courseSlot.CourseSlots, error) {
	out := make(courseSlot.CourseSlots, 0, len(slotIDs))
	for _, id := range slotIDs {
		slot, err := c.MustGet(ctx, id)
		if err != nil {
			return nil, err
		}
		out = append(out, slot)
	}
	return out, nil
}

// GetTeacherSlots 取该讲师现有的全部排期。
func (c *CourseSlotImp) GetTeacherSlots(ctx context.Context, teacherID int64) (courseSlot.CourseSlots, error) {
	rows, err := c.data.Conn(ctx).QueryContext(ctx, `
SELECT `+help.CourseSlotColumns+`
  FROM course_slot
 WHERE teacher_id = ?
 ORDER BY id`, teacherID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	return scanSlots(rows)
}

// GetClassroomSlots 取该教室现有的全部排期。
func (c *CourseSlotImp) GetClassroomSlots(ctx context.Context, classroomID string) (courseSlot.CourseSlots, error) {
	rows, err := c.data.Conn(ctx).QueryContext(ctx, `
SELECT `+help.CourseSlotColumns+`
  FROM course_slot
 WHERE classroom_id = ?
 ORDER BY id`, classroomID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	return scanSlots(rows)
}

// GetCourseSlots 取该课程现有的全部排期。
func (c *CourseSlotImp) GetCourseSlots(ctx context.Context, courseID string) (courseSlot.CourseSlots, error) {
	rows, err := c.data.Conn(ctx).QueryContext(ctx, `
SELECT `+help.CourseSlotColumns+`
  FROM course_slot
 WHERE course_id = ?
 ORDER BY id`, courseID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	return scanSlots(rows)
}

// GetStudentSlots 取该学员在学课程的全部排期。
func (c *CourseSlotImp) GetStudentSlots(ctx context.Context, studentID int64) (courseSlot.CourseSlots, error) {
	rows, err := c.data.Conn(ctx).QueryContext(ctx, `
SELECT `+help.CourseSlotColumns+`
  FROM course_slot
 WHERE course_id IN (
       SELECT course_id FROM course_enrollment
        WHERE student_id = ? AND status = 1 -- 1 = 在读
 )
 ORDER BY id`, studentID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	return scanSlots(rows)
}

// scanSlots 把多行扫成 CourseSlots。
func scanSlots(rows *sql.Rows) (courseSlot.CourseSlots, error) {
	out := make(courseSlot.CourseSlots, 0)
	for rows.Next() {
		po, err := help.ScanCourseSlot(rows)
		if err != nil {
			return nil, err
		}
		do, err := model.CourseSlotToDO(po)
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
