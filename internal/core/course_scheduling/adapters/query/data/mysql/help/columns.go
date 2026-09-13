package help

import "strings"

// 每张表的列清单。顺序必须与建表脚本、与 model 包里字段的声明顺序一致 ——
// ScanXxx 是按位置扫的，这里错了不会报错，只会静静地把值塞进错的字段。
//
// 刻意显式列出列名而不用 SELECT *：表结构一变，这里立刻对不上。
const (
	CourseTypeColumns       = `id, name, description, created_at, updated_at, lock_version`
	ClassroomColumns        = `id, building, floor, room, capacity, allocated, status, lock_version`
	StudentColumns          = `id, name, student_type, phone, email, status, created_at, updated_at, lock_version`
	TeacherColumns          = `id, student_id, name, title, phone, email, status, created_at, updated_at, lock_version`
	CourseColumns           = `id, course_type_id, capacity_max, capacity_enrolled, enroll_start_at, enroll_end_at, drop_deadline, period_start_at, period_end_at, total_hours, completed_hours, lock_version`
	CourseSlotColumns       = `id, course_id, weekday, start_time, end_time, teacher_id, classroom_id, created_at, updated_at, lock_version`
	CourseSlotChangeColumns = `id, course_id, applicant_id, change_type, original_slot_id, original_date, original_teacher_id, original_classroom_id, original_start_time, original_end_time, target_start_at, target_end_at, target_teacher_id, target_classroom_id, reason, created_at, updated_at, lock_version`
	EnrollmentColumns       = `id, student_id, course_id, status, enrolled_at, completed_at, dropped_at, updated_at, lock_version`
	AbsenceColumns          = `id, student_id, course_id, course_slot_id, schedule_date, missed_hours, absence_type, reason, created_at, updated_at, lock_version`
	MakeupColumns           = `id, student_id, course_id, original_slot_id, original_date, target_slot_id, target_date, makeup_hours, status, completed_at, created_at, updated_at, lock_version`
	QualificationColumns    = `id, teacher_id, course_type_id, certified_at, status, updated_at, lock_version`
)

// Qualify 给列清单加上表别名前缀，用于带 JOIN / 子查询的语句。
//
// 例如 Qualify(CourseColumns, "c") -> `c.id, c.course_type_id, …`，
// 这样外层 SELECT 里不会和子查询的同名列撞上。
func Qualify(columns, alias string) string {
	parts := strings.Split(columns, ", ")
	for i, part := range parts {
		parts[i] = alias + "." + part
	}
	return strings.Join(parts, ", ")
}
