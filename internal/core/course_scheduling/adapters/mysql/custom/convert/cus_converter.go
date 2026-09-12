// Package custom 手写实现 mysql 包定义的 DO <-> PO 转换契约。
//
// 整个包就一个类型：Converter，它实现了 mysql.Converter 的全部 22 个方法
// （11 个聚合 × 进出两个方向）。文件末尾用编译期断言把接口和实现对起来，
// 以后接口改签名或漏实现方法，go build 就会报错。
//
// 读路径（XxxToDO）刻意走一遍领域模型的值对象构造函数（NewLocation / NewCapacity /
// NewCoursePeriod / NewDayTimeRange / NewTargetPlan …），让非法数据在这一层
// 就被拦下来，而不是悄悄造出一个不满足不变量的聚合。
//
// 写路径（XxxToPO）只取值打平，不写 LockVersion：乐观锁是存储层的事，
// 由 repo 在 UPDATE / DELETE 语句里带上版本号。
package custom

import (
	"database/sql"
	"fmt"
	"time"

	"cqrs/internal/core/course_scheduling/adapters/mysql"
	"cqrs/internal/core/course_scheduling/adapters/mysql/custom/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
)

// Converter 手写的转换器：mysql.Converter 的唯一实现。
//
// 无状态，零值可用，wire 里直接 new(custom.Converter) 就行。
type Converter struct{}

// 编译期断言：实现必须与接口同步。
var _ mysql.Converter = Converter{}

// --- course_type ---

// CourseTypeToPO 写路径。
func (Converter) CourseTypeToPO(do *courseType.CourseType) *model.CourseType {
	return &model.CourseType{
		ID:          do.ID(),
		Name:        do.Name(),
		Description: do.Description(),
		CreatedAt:   do.CreatedAt(),
		UpdatedAt:   do.UpdatedAt(),
	}
}

// CourseTypeToDO 读路径。
func (Converter) CourseTypeToDO(po *model.CourseType) (*courseType.CourseType, error) {
	return courseType.Reconstitute(
		po.ID, po.Name, po.Description, po.CreatedAt, po.UpdatedAt,
	), nil
}

// --- classroom ---

// ClassroomToPO 写路径：Location 值对象打平成三列。
func (Converter) ClassroomToPO(do *classroom.Classroom) *model.Classroom {
	location := do.Location()
	return &model.Classroom{
		ID:        do.ID(),
		Building:  location.Building(),
		Floor:     location.Floor(),
		Room:      location.Room(),
		Capacity:  do.Capacity(),
		Allocated: do.AllocatedSeats(),
		Status:    uint8(do.Status()),
	}
}

// ClassroomToDO 读路径：三列重新拼回 Location，楼栋/房间号为空会被拒绝。
func (Converter) ClassroomToDO(po *model.Classroom) (*classroom.Classroom, error) {
	location, err := classroom.NewLocation(po.Building, po.Floor, po.Room)
	if err != nil {
		return nil, err
	}
	return classroom.Reconstitute(
		po.ID, location, po.Capacity, po.Allocated, classroom.Status(po.Status),
	), nil
}

// --- student ---

// StudentToPO 写路径：ContactInfo 值对象打平成 phone / email。
func (Converter) StudentToPO(do *student.Student) *model.Student {
	contact := do.Contact()
	return &model.Student{
		ID:          do.ID(),
		Name:        do.Name(),
		StudentType: uint8(do.StudentType()),
		Phone:       contact.Phone(),
		Email:       contact.Email(),
		Status:      uint8(do.Status()),
		CreatedAt:   do.CreatedAt(),
		UpdatedAt:   do.UpdatedAt(),
	}
}

// StudentToDO 读路径：联系方式要过一遍格式校验（手机号必填）。
func (Converter) StudentToDO(po *model.Student) (*student.Student, error) {
	contact, err := student.NewContactInfo(po.Phone, po.Email)
	if err != nil {
		return nil, err
	}
	return student.Reconstitute(
		po.ID, po.Name, student.StudentType(po.StudentType), contact,
		student.Status(po.Status), po.CreatedAt, po.UpdatedAt,
	), nil
}

// --- teacher ---

// TeacherToPO 写路径：ContactInfo 值对象打平成 phone / email。
func (Converter) TeacherToPO(do *teacher.Teacher) *model.Teacher {
	contact := do.Contact()
	return &model.Teacher{
		ID:        do.ID(),
		StudentID: do.StudentID(),
		Name:      do.Name(),
		Title:     do.Title(),
		Phone:     contact.Phone(),
		Email:     contact.Email(),
		Status:    uint8(do.Status()),
		CreatedAt: do.CreatedAt(),
		UpdatedAt: do.UpdatedAt(),
	}
}

// TeacherToDO 读路径：联系方式要过一遍格式校验（手机号必填）。
func (Converter) TeacherToDO(po *model.Teacher) (*teacher.Teacher, error) {
	contact, err := teacher.NewContactInfo(po.Phone, po.Email)
	if err != nil {
		return nil, err
	}
	return teacher.Reconstitute(
		po.ID, po.StudentID, po.Name, po.Title, contact,
		teacher.Status(po.Status), po.CreatedAt, po.UpdatedAt,
	), nil
}

// --- course ---

// CourseToPO 写路径：Capacity / EnrollmentWindow / CoursePeriod 三个值对象打平。
func (Converter) CourseToPO(do *course.Course) *model.Course {
	capacity := do.Capacity()
	window := do.Enrollment()
	period := do.Period()
	return &model.Course{
		ID:               do.ID(),
		CourseTypeID:     do.CourseTypeID(),
		CapacityMax:      capacity.Max(),
		CapacityEnrolled: capacity.Enrolled(),
		EnrollStartAt:    window.StartAt(),
		EnrollEndAt:      window.EndAt(),
		DropDeadline:     window.DropDeadline(),
		PeriodStartAt:    period.StartAt(),
		PeriodEndAt:      period.EndAt(),
		TotalHours:       period.TotalHours(),
		CompletedHours:   period.CompletedHours(),
	}
}

// CourseToDO 读路径：容量与课时都走值对象构造函数，越界数据在这里报错。
func (Converter) CourseToDO(po *model.Course) (*course.Course, error) {
	capacity, err := course.NewCapacity(po.CapacityMax, po.CapacityEnrolled)
	if err != nil {
		return nil, err
	}
	period, err := course.NewCoursePeriod(
		po.PeriodStartAt, po.PeriodEndAt, po.TotalHours, po.CompletedHours,
	)
	if err != nil {
		return nil, err
	}
	window := course.NewEnrollmentWindow(po.EnrollStartAt, po.EnrollEndAt, po.DropDeadline)

	return course.Reconstitute(po.ID, po.CourseTypeID, capacity, window, period), nil
}

// --- course_slot ---

// CourseSlotToPO 写路径：DayTimeRange 打平成两个 time 列。
func (Converter) CourseSlotToPO(do *courseSlot.CourseSlot) *model.CourseSlot {
	timeRange := do.TimeRange()
	return &model.CourseSlot{
		ID:          do.ID(),
		CourseID:    do.CourseID(),
		Weekday:     uint8(do.Weekday()),
		StartTime:   dayTimeToPO(timeRange.Start()),
		EndTime:     dayTimeToPO(timeRange.End()),
		TeacherID:   do.TeacherID(),
		ClassroomID: do.ClassroomID(),
		CreatedAt:   do.CreatedAt(),
		UpdatedAt:   do.UpdatedAt(),
	}
}

// CourseSlotToDO 读路径：两个 time 列拼回 DayTimeRange，先后顺序反了会被拒绝。
func (Converter) CourseSlotToDO(po *model.CourseSlot) (*courseSlot.CourseSlot, error) {
	start, err := dayTimeFromPO(po.StartTime)
	if err != nil {
		return nil, err
	}
	end, err := dayTimeFromPO(po.EndTime)
	if err != nil {
		return nil, err
	}
	timeRange, err := courseSlot.NewDayTimeRange(start, end)
	if err != nil {
		return nil, err
	}
	return courseSlot.Reconstitute(
		po.ID, po.CourseID, time.Weekday(po.Weekday), timeRange,
		po.TeacherID, po.ClassroomID, po.CreatedAt, po.UpdatedAt,
	), nil
}

// --- course_slot_change ---

// CourseSlotChangeToPO 写路径：OriginalPlan 打平成 original_* 六列，
// TargetPlan 打平成 target_* 四列。
func (Converter) CourseSlotChangeToPO(do *courseSlotChange.CourseSlotChange) *model.CourseSlotChange {
	original := do.OriginalPlan()
	target := do.TargetPlan()
	return &model.CourseSlotChange{
		ID:          do.ID(),
		CourseID:    do.CourseID(),
		ApplicantID: do.ApplicantID(),
		ChangeType:  uint8(do.ChangeType()),

		OriginalSlotID:      original.SlotID(),
		OriginalDate:        original.Date(),
		OriginalTeacherID:   original.TeacherID(),
		OriginalClassroomID: original.ClassroomID(),
		OriginalStartTime:   original.StartTimeStr(),
		OriginalEndTime:     original.EndTimeStr(),

		TargetStartAt:     target.TargetStartAt(),
		TargetEndAt:       target.TargetEndAt(),
		TargetTeacherID:   target.TeacherID(),
		TargetClassroomID: target.ClassroomID(),

		Reason:    do.Reason(),
		CreatedAt: do.CreatedAt(),
		UpdatedAt: do.UpdatedAt(),
	}
}

// CourseSlotChangeToDO 读路径：目标计划要走 NewTargetPlan，跨天或时间倒挂会报错。
//
// 原计划的 original_start_time / original_end_time 表里就是字符串（varchar(5)），
// 和值对象一一对应，直接带过去。
func (Converter) CourseSlotChangeToDO(po *model.CourseSlotChange) (*courseSlotChange.CourseSlotChange, error) {
	original := courseSlotChange.NewOriginalPlan(
		po.OriginalSlotID, po.OriginalDate, po.OriginalTeacherID,
		po.OriginalClassroomID, po.OriginalStartTime, po.OriginalEndTime,
	)
	target, err := courseSlotChange.NewTargetPlan(
		po.TargetStartAt, po.TargetEndAt, po.TargetTeacherID, po.TargetClassroomID,
	)
	if err != nil {
		return nil, err
	}
	return courseSlotChange.Reconstitute(
		po.ID, po.CourseID, po.ApplicantID, courseSlotChange.ChangeType(po.ChangeType),
		original, target, po.Reason, po.CreatedAt, po.UpdatedAt,
	), nil
}

// --- course_enrollment ---

// EnrollmentToPO 写路径：零值时间表示「没有」，写进可空列就是 NULL。
func (Converter) EnrollmentToPO(do *enrollment.CourseEnrollment) *model.CourseEnrollment {
	return &model.CourseEnrollment{
		ID:          do.ID(),
		StudentID:   do.StudentID(),
		CourseID:    do.CourseID(),
		Status:      uint8(do.Status()),
		EnrolledAt:  do.EnrolledAt(),
		CompletedAt: nullTime(do.CompletedAt()),
		DroppedAt:   nullTime(do.DroppedAt()),
		UpdatedAt:   do.UpdatedAt(),
	}
}

// EnrollmentToDO 读路径：NULL 还原成零值时间。
func (Converter) EnrollmentToDO(po *model.CourseEnrollment) (*enrollment.CourseEnrollment, error) {
	return enrollment.Reconstitute(
		po.ID, po.StudentID, po.CourseID, enrollment.Status(po.Status),
		po.EnrolledAt, timeValue(po.CompletedAt), timeValue(po.DroppedAt), po.UpdatedAt,
	), nil
}

// --- absence_record ---

// AbsenceToPO 写路径。这张表没有状态列，缺勤就是一条事实。
func (Converter) AbsenceToPO(do *absence.AbsenceRecord) *model.AbsenceRecord {
	return &model.AbsenceRecord{
		ID:           do.ID(),
		StudentID:    do.StudentID(),
		CourseID:     do.CourseID(),
		CourseSlotID: do.CourseSlotID(),
		ScheduleDate: do.ScheduleDate(),
		MissedHours:  do.MissedHours(),
		AbsenceType:  uint8(do.AbsenceType()),
		Reason:       do.Reason(),
		CreatedAt:    do.CreatedAt(),
		UpdatedAt:    do.UpdatedAt(),
	}
}

// AbsenceToDO 读路径。
func (Converter) AbsenceToDO(po *model.AbsenceRecord) (*absence.AbsenceRecord, error) {
	return absence.Reconstitute(
		po.ID, po.StudentID, po.CourseID, po.CourseSlotID, po.ScheduleDate,
		po.MissedHours, absence.AbsenceType(po.AbsenceType), po.Reason,
		po.CreatedAt, po.UpdatedAt,
	), nil
}

// --- student_makeup ---

// MakeupToPO 写路径。
func (Converter) MakeupToPO(do *makeup.StudentMakeup) *model.StudentMakeup {
	return &model.StudentMakeup{
		ID:             do.ID(),
		StudentID:      do.StudentID(),
		CourseID:       do.CourseID(),
		OriginalSlotID: do.OriginalSlotID(),
		OriginalDate:   do.OriginalDate(),
		TargetSlotID:   do.TargetSlotID(),
		TargetDate:     do.TargetDate(),
		MakeupHours:    do.MakeupHours(),
		Status:         uint8(do.Status()),
		CompletedAt:    nullTime(do.CompletedAt()),
		CreatedAt:      do.CreatedAt(),
		UpdatedAt:      do.UpdatedAt(),
	}
}

// MakeupToDO 读路径：completed_at 为 NULL 时还原成零值时间（= 还没核销）。
func (Converter) MakeupToDO(po *model.StudentMakeup) (*makeup.StudentMakeup, error) {
	return makeup.Reconstitute(
		po.ID, po.StudentID, po.CourseID, po.OriginalSlotID, po.OriginalDate,
		po.TargetSlotID, po.TargetDate, po.MakeupHours, makeup.Status(po.Status),
		timeValue(po.CompletedAt), po.CreatedAt, po.UpdatedAt,
	), nil
}

// --- qualification ---

// QualificationToPO 写路径。这张表没有 created_at，认证时间就是 certified_at。
func (Converter) QualificationToPO(do *qualification.Qualification) *model.Qualification {
	return &model.Qualification{
		ID:           do.ID(),
		TeacherID:    do.TeacherID(),
		CourseTypeID: do.CourseTypeID(),
		CertifiedAt:  do.CertifiedAt(),
		Status:       uint8(do.Status()),
		UpdatedAt:    do.UpdatedAt(),
	}
}

// QualificationToDO 读路径。
func (Converter) QualificationToDO(po *model.Qualification) (*qualification.Qualification, error) {
	return qualification.Reconstitute(
		po.ID, po.TeacherID, po.CourseTypeID, po.CertifiedAt,
		qualification.Status(po.Status), po.UpdatedAt,
	), nil
}

// --- 小工具 ---

// nullTime 领域模型的零值时间 -> 可空列（零值 = NULL）。
func nullTime(t time.Time) sql.NullTime {
	return sql.NullTime{Time: t, Valid: !t.IsZero()}
}

// timeValue 可空列 -> 领域模型的零值时间（NULL = 零值）。
func timeValue(n sql.NullTime) time.Time {
	if !n.Valid {
		return time.Time{}
	}
	return n.Time
}

// dayTimeToPO DayTime 值对象 -> time 列，格式固定成 "15:04:05"。
//
// 不直接用 DayTime.String()（它给的是 "15:04"）—— 让写入的字符串
// 和驱动读出来的格式一致，日志和对账时省得两副面孔。
func dayTimeToPO(t courseSlot.DayTime) string {
	minutes := t.TotalMinutes()
	return fmt.Sprintf("%02d:%02d:00", minutes/60, minutes%60)
}

// dayTimeFromPO time 列 -> DayTime 值对象。
//
// 驱动给的是 "15:04:05"；这里顺手兼容 "15:04"，两种都认。
func dayTimeFromPO(s string) (courseSlot.DayTime, error) {
	for _, layout := range []string{"15:04:05", "15:04"} {
		parsed, err := time.Parse(layout, s)
		if err == nil {
			return courseSlot.NewDayTime(parsed.Hour(), parsed.Minute())
		}
	}
	return courseSlot.DayTime{}, fmt.Errorf("model: 无法解析时间列 %q", s)
}
