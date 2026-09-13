package help

import "cqrs/internal/core/course_scheduling/adapters/command/data/mysql/model"

// Scanner 抽象 *sql.Rows 与 *sql.Row 共有的 Scan 能力，让同一段扫行代码
// 既能用于 Query（多行）也能用于 QueryRow（单行）。
type Scanner interface {
	Scan(dest ...any) error
}

func ScanClassroom(sc Scanner) (*model.Classroom, error) {
	po := &model.Classroom{}
	if err := sc.Scan(
		&po.ID, &po.Building, &po.Floor, &po.Room,
		&po.Capacity, &po.Allocated, &po.Status,
	); err != nil {
		return nil, err
	}
	return po, nil
}

func ScanStudent(sc Scanner) (*model.Student, error) {
	po := &model.Student{}
	if err := sc.Scan(
		&po.ID, &po.Name, &po.StudentType, &po.Phone, &po.Email,
		&po.Status, &po.CreatedAt, &po.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return po, nil
}

func ScanTeacher(sc Scanner) (*model.Teacher, error) {
	po := &model.Teacher{}
	if err := sc.Scan(
		&po.ID, &po.StudentID, &po.Name, &po.Title, &po.Phone, &po.Email,
		&po.Status, &po.CreatedAt, &po.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return po, nil
}

func ScanCourse(sc Scanner) (*model.Course, error) {
	po := &model.Course{}
	if err := sc.Scan(
		&po.ID, &po.CourseTypeID, &po.CapacityMax, &po.CapacityEnrolled,
		&po.EnrollStartAt, &po.EnrollEndAt, &po.DropDeadline,
		&po.PeriodStartAt, &po.PeriodEndAt, &po.TotalHours, &po.CompletedHours,
	); err != nil {
		return nil, err
	}
	return po, nil
}

func ScanCourseSlot(sc Scanner) (*model.CourseSlot, error) {
	po := &model.CourseSlot{}
	if err := sc.Scan(
		&po.ID, &po.CourseID, &po.Weekday, &po.StartTime, &po.EndTime,
		&po.TeacherID, &po.ClassroomID, &po.CreatedAt, &po.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return po, nil
}

func ScanCourseSlotChange(sc Scanner) (*model.CourseSlotChange, error) {
	po := &model.CourseSlotChange{}
	if err := sc.Scan(
		&po.ID, &po.CourseID, &po.ApplicantID, &po.ChangeType,
		&po.OriginalSlotID, &po.OriginalDate, &po.OriginalTeacherID, &po.OriginalClassroomID,
		&po.OriginalStartTime, &po.OriginalEndTime,
		&po.TargetStartAt, &po.TargetEndAt, &po.TargetTeacherID, &po.TargetClassroomID,
		&po.Reason, &po.CreatedAt, &po.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return po, nil
}

func ScanCourseType(sc Scanner) (*model.CourseType, error) {
	po := &model.CourseType{}
	if err := sc.Scan(
		&po.ID, &po.Name, &po.Description, &po.CreatedAt, &po.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return po, nil
}

func ScanEnrollment(sc Scanner) (*model.CourseEnrollment, error) {
	po := &model.CourseEnrollment{}
	if err := sc.Scan(
		&po.ID, &po.StudentID, &po.CourseID, &po.Status, &po.EnrolledAt,
		&po.CompletedAt, &po.DroppedAt, &po.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return po, nil
}

func ScanAbsence(sc Scanner) (*model.AbsenceRecord, error) {
	po := &model.AbsenceRecord{}
	if err := sc.Scan(
		&po.ID, &po.StudentID, &po.CourseID, &po.CourseSlotID, &po.ScheduleDate,
		&po.MissedHours, &po.AbsenceType, &po.Reason, &po.CreatedAt, &po.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return po, nil
}

func ScanMakeup(sc Scanner) (*model.StudentMakeup, error) {
	po := &model.StudentMakeup{}
	if err := sc.Scan(
		&po.ID, &po.StudentID, &po.CourseID,
		&po.OriginalSlotID, &po.OriginalDate, &po.TargetSlotID, &po.TargetDate,
		&po.MakeupHours, &po.Status, &po.CompletedAt,
		&po.CreatedAt, &po.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return po, nil
}

func ScanQualification(sc Scanner) (*model.Qualification, error) {
	po := &model.Qualification{}
	if err := sc.Scan(
		&po.ID, &po.TeacherID, &po.CourseTypeID, &po.CertifiedAt,
		&po.Status, &po.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return po, nil
}
