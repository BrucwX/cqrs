package imp

import (
	"context"
	"database/sql"

	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/adapters/command/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type AssignTeacherImp struct {
	data *mysql.Data
}

// 编译期断言：实现必须满足接口。
var _ repo.AssignTeacherRepo = (*AssignTeacherImp)(nil)

// NewClassroomQuery 创建 MySQL 版教室查询。
func NewAssignTeacherImp(d *mysql.Data) repo.AssignTeacherRepo {
	return &AssignTeacherImp{data: d}
}

func (c *AssignTeacherImp) GetSlots(slotIDs []string) (courseSlot.CourseSlots, error) {
	panic("implement me")
}

func (c *AssignTeacherImp) GetTeacher(teacherID int64) (teacher.Teacher, error) {
	panic("implement me")
}

func (c *AssignTeacherImp) GetTeacherSlots(teacherID int64) (courseSlot.CourseSlots, error) {
	ctx := context.Background()
	const query = `SELECT id, course_id, weekday, start_time, end_time, teacher_id, classroom_id, created_at, updated_at, lock_version
		FROM course_slot WHERE teacher_id = ?`

	rows, err := c.data.Conn(ctx).QueryContext(ctx, query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots courseSlot.CourseSlots
	for rows.Next() {
		var po model.CourseSlot
		if err := rows.Scan(
			&po.ID, &po.CourseID, &po.Weekday, &po.StartTime, &po.EndTime,
			&po.TeacherID, &po.ClassroomID, &po.CreatedAt, &po.UpdatedAt, &po.LockVersion,
		); err != nil {
			return nil, err
		}
		do, err := model.CourseSlotToDO(&po)
		if err != nil {
			return nil, err
		}
		slots = append(slots, *do)
	}
	return slots, rows.Err()
}

func (c *AssignTeacherImp) GetQualifications(teacherID int64) ([]qualification.Qualification, error) {
	ctx := context.Background()
	const query = `SELECT id, teacher_id, course_type_id, certified_at, status, updated_at, lock_version
		FROM qualification WHERE teacher_id = ?`

	rows, err := c.data.Conn(ctx).QueryContext(ctx, query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quals []qualification.Qualification
	for rows.Next() {
		var po model.Qualification
		if err := rows.Scan(
			&po.ID, &po.TeacherID, &po.CourseTypeID, &po.CertifiedAt,
			&po.Status, &po.UpdatedAt, &po.LockVersion,
		); err != nil {
			return nil, err
		}
		do, err := model.QualificationToDO(&po)
		if err != nil {
			return nil, err
		}
		quals = append(quals, *do)
	}
	return quals, rows.Err()
}

func (c *AssignTeacherImp) GetCourseType(courseID string) (courseType.CourseType, error) {
	ctx := context.Background()
	const query = `SELECT ct.id, ct.name, ct.description, ct.created_at, ct.updated_at, ct.lock_version
		FROM course_type ct
		JOIN course c ON c.course_type_id = ct.id
		WHERE c.id = ?`

	var po model.CourseType
	err := c.data.Conn(ctx).QueryRowContext(ctx, query, courseID).Scan(
		&po.ID, &po.Name, &po.Description, &po.CreatedAt, &po.UpdatedAt, &po.LockVersion,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return courseType.CourseType{}, repo.ErrCourseTypeNotFound
		}
		return courseType.CourseType{}, err
	}
	do, err := model.CourseTypeToDO(&po)
	if err != nil {
		return courseType.CourseType{}, err
	}
	return *do, nil
}
