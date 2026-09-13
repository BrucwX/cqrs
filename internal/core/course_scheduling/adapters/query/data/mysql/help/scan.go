// Package help 是查询侧 MySQL 适配器的辅助函数：把 database/sql 的一行
// 扫成 model 里的 PO，把 PO 交给转换函数变成 DO，顺带分页换算。
//
// 与 implement 分开的原因：列清单和 Scan 参数是「表结构」的事，不是「查询语义」
// 的事。implement 里只留每个接口一份 SQL，这层要换（换驱动、加埋点、加缓存）
// 不用动任何一个实现。
//
// 四个文件：
//
//	columns.go  每张表的列清单 + Qualify（给列加表别名）
//	scan.go     一行 -> 一个 PO：Scanner 接口 + 11 个 ScanXxx
//	exec.go     执行与转换：List / QueryAll / Exists
//	page.go     分页换算：LimitOffset
//
// 最省事的用法是 QueryAll，一次把「执行 SQL + 扫行 + 转 DO」走完：
//
//	pos, err := help.QueryAll(ctx, q,
//		`SELECT `+help.CourseColumns+` FROM course`,
//		nil, help.ScanCourse, model.CourseToDO)
package help

import (
	"cqrs/internal/core/course_scheduling/adapters/query/data/mysql/model"
)

// Scanner 抽象 *sql.Rows 与 *sql.Row 共有的 Scan 能力，让同一段扫行代码
// 既能用于 Query（多行）也能用于 QueryRow（单行）。
type Scanner interface {
	Scan(dest ...any) error
}

func ScanCourseType(sc Scanner) (*model.CourseType, error) {
	var po model.CourseType
	if err := sc.Scan(
		&po.ID, &po.Name, &po.Description, &po.CreatedAt, &po.UpdatedAt, &po.LockVersion,
	); err != nil {
		return nil, err
	}
	return &po, nil
}

func ScanClassroom(sc Scanner) (*model.Classroom, error) {
	var po model.Classroom
	if err := sc.Scan(
		&po.ID, &po.Building, &po.Floor, &po.Room,
		&po.Capacity, &po.Allocated, &po.Status, &po.LockVersion,
	); err != nil {
		return nil, err
	}
	return &po, nil
}

func ScanStudent(sc Scanner) (*model.Student, error) {
	var po model.Student
	if err := sc.Scan(
		&po.ID, &po.Name, &po.StudentType, &po.Phone, &po.Email,
		&po.Status, &po.CreatedAt, &po.UpdatedAt, &po.LockVersion,
	); err != nil {
		return nil, err
	}
	return &po, nil
}

func ScanTeacher(sc Scanner) (*model.Teacher, error) {
	var po model.Teacher
	if err := sc.Scan(
		&po.ID, &po.StudentID, &po.Name, &po.Title, &po.Phone, &po.Email,
		&po.Status, &po.CreatedAt, &po.UpdatedAt, &po.LockVersion,
	); err != nil {
		return nil, err
	}
	return &po, nil
}

func ScanCourse(sc Scanner) (*model.Course, error) {
	var po model.Course
	if err := sc.Scan(
		&po.ID, &po.CourseTypeID, &po.CapacityMax, &po.CapacityEnrolled,
		&po.EnrollStartAt, &po.EnrollEndAt, &po.DropDeadline,
		&po.PeriodStartAt, &po.PeriodEndAt, &po.TotalHours, &po.CompletedHours,
		&po.LockVersion,
	); err != nil {
		return nil, err
	}
	return &po, nil
}

func ScanCourseSlot(sc Scanner) (*model.CourseSlot, error) {
	var po model.CourseSlot
	if err := sc.Scan(
		&po.ID, &po.CourseID, &po.Weekday, &po.StartTime, &po.EndTime,
		&po.TeacherID, &po.ClassroomID, &po.CreatedAt, &po.UpdatedAt, &po.LockVersion,
	); err != nil {
		return nil, err
	}
	return &po, nil
}

func ScanCourseSlotChange(sc Scanner) (*model.CourseSlotChange, error) {
	var po model.CourseSlotChange
	if err := sc.Scan(
		&po.ID, &po.CourseID, &po.ApplicantID, &po.ChangeType,
		&po.OriginalSlotID, &po.OriginalDate, &po.OriginalTeacherID, &po.OriginalClassroomID,
		&po.OriginalStartTime, &po.OriginalEndTime,
		&po.TargetStartAt, &po.TargetEndAt, &po.TargetTeacherID, &po.TargetClassroomID,
		&po.Reason, &po.CreatedAt, &po.UpdatedAt, &po.LockVersion,
	); err != nil {
		return nil, err
	}
	return &po, nil
}

func ScanEnrollment(sc Scanner) (*model.CourseEnrollment, error) {
	var po model.CourseEnrollment
	if err := sc.Scan(
		&po.ID, &po.StudentID, &po.CourseID, &po.Status, &po.EnrolledAt,
		&po.CompletedAt, &po.DroppedAt, &po.UpdatedAt, &po.LockVersion,
	); err != nil {
		return nil, err
	}
	return &po, nil
}

func ScanAbsence(sc Scanner) (*model.AbsenceRecord, error) {
	var po model.AbsenceRecord
	if err := sc.Scan(
		&po.ID, &po.StudentID, &po.CourseID, &po.CourseSlotID, &po.ScheduleDate,
		&po.MissedHours, &po.AbsenceType, &po.Reason,
		&po.CreatedAt, &po.UpdatedAt, &po.LockVersion,
	); err != nil {
		return nil, err
	}
	return &po, nil
}

func ScanMakeup(sc Scanner) (*model.StudentMakeup, error) {
	var po model.StudentMakeup
	if err := sc.Scan(
		&po.ID, &po.StudentID, &po.CourseID,
		&po.OriginalSlotID, &po.OriginalDate, &po.TargetSlotID, &po.TargetDate,
		&po.MakeupHours, &po.Status, &po.CompletedAt,
		&po.CreatedAt, &po.UpdatedAt, &po.LockVersion,
	); err != nil {
		return nil, err
	}
	return &po, nil
}

func ScanQualification(sc Scanner) (*model.Qualification, error) {
	var po model.Qualification
	if err := sc.Scan(
		&po.ID, &po.TeacherID, &po.CourseTypeID, &po.CertifiedAt,
		&po.Status, &po.UpdatedAt, &po.LockVersion,
	); err != nil {
		return nil, err
	}
	return &po, nil
}
