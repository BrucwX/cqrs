package help

// Package help 是写侧 MySQL 适配器的辅助函数：把 database/sql 的一行
// 扫成 model 里的 PO。
//
// 与读侧（adapters/query/data/mysql/implement/help）一一对应：列清单与扫行是「表结构」
// 的事，implement 里只留每个仓库一份 SQL。
//
// ⚠️ 目前只有 classroom 启用了乐观锁（列清单带 lock_version、写回时校验版本）；
// 其余表的 lock_version 列建好了但还没读，接的时候照 classroom 的写法来。

// 每张表的列清单。顺序必须与建表脚本、与 model 包里字段声明顺序一致 ——
// ScanXxx 是按位置扫的，这里错了不会报错，只会静静地把值塞进错的字段。
// 带 lock_version 的表要把它放在**最后**（跟 model 的字段顺序对齐）。
const (
	ClassroomColumns        = `id, building, floor, room, capacity, allocated, status, lock_version`
	StudentColumns          = `id, name, student_type, phone, email, status, created_at, updated_at`
	TeacherColumns          = `id, student_id, name, title, phone, email, status, created_at, updated_at`
	CourseColumns           = `id, course_type_id, capacity_max, capacity_enrolled, enroll_start_at, enroll_end_at, drop_deadline, period_start_at, period_end_at, total_hours, completed_hours`
	CourseSlotColumns       = `id, course_id, weekday, start_time, end_time, teacher_id, classroom_id, created_at, updated_at`
	CourseSlotChangeColumns = `id, course_id, applicant_id, change_type, original_slot_id, original_date, original_teacher_id, original_classroom_id, original_start_time, original_end_time, target_start_at, target_end_at, target_teacher_id, target_classroom_id, reason, created_at, updated_at`
	CourseTypeColumns       = `id, name, description, created_at, updated_at`
	EnrollmentColumns       = `id, student_id, course_id, status, enrolled_at, completed_at, dropped_at, updated_at`
	AbsenceColumns          = `id, student_id, course_id, course_slot_id, schedule_date, missed_hours, absence_type, reason, created_at, updated_at`
	MakeupColumns           = `id, student_id, course_id, original_slot_id, original_date, target_slot_id, target_date, makeup_hours, status, completed_at, created_at, updated_at`
	QualificationColumns    = `id, teacher_id, course_type_id, certified_at, status, updated_at`
)
