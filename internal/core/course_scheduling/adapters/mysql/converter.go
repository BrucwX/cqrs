package mysql

import (
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

// 本文件定义数据模型（PO，custom/model）与领域模型（DO，domain/aggregate）
// 之间的转换契约：一个 Converter 接口，11 个聚合的往返转换都收在里面。
//
// 为什么接口留在这里、实现放到 custom 包：
//
//   - repo（command / query）只依赖接口，不依赖具体实现，将来换一套转换逻辑
//     （比如换成生成器产物）不用动 repo；
//   - 接口是「调用方需要什么」，实现是「具体怎么转」，分开也好逐个对照检查。
//
// 方法命名：<聚合>ToPO / <聚合>ToDO。
//
//	不是 ToPO / ToDO —— Go 的接口里不允许出现两个同名不同参的方法，
//	11 个聚合挤在一个接口里就必须带聚合前缀。
//
// 方向约定：
//
//   - <聚合>ToPO 是写路径，DO -> PO。纯取值打平，不会失败，所以不返回 error。
//   - <聚合>ToDO 是读路径，PO -> DO。领域模型的值对象是自校验的（手机号格式、
//     容量上下界、时间段先后、目标计划不能跨天…），脏数据得在这里被拦下来，
//     所以返回 error。
//   - ToDO 一律返回 error，即使某个聚合的 Reconstitute 现在不会失败 ——
//     保持 22 个方法形状一致，读路径的调用方写法统一。
//
// LockVersion：这是存储层的乐观锁字段，领域模型里没有对应概念。
// 所以 ToPO 不给它赋值（保持 0），由 repo 在写语句里显式带上版本号；
// ToDO 也不读它，领域模型不感知乐观锁。
//
// 值对象一律就地打平，不另开类型：
//
//	Location         -> building / floor / room
//	Capacity         -> capacity_max / capacity_enrolled
//	EnrollmentWindow -> enroll_start_at / enroll_end_at / drop_deadline
//	CoursePeriod     -> period_start_at / period_end_at / total_hours / completed_hours
//	ContactInfo      -> phone / email
//	DayTimeRange     -> start_time / end_time
//	OriginalPlan     -> original_* 六列
//	TargetPlan       -> target_* 四列
//
// 实现见 custom 包（手写，不用代码生成器）：custom.Converter。
type Converter interface {
	// --- course_type <-> courseType.CourseType ---

	// CourseTypeToPO 领域模型 -> 数据模型（写路径）。
	CourseTypeToPO(do *courseType.CourseType) *model.CourseType
	// CourseTypeToDO 数据模型 -> 领域模型（读路径）。
	CourseTypeToDO(po *model.CourseType) (*courseType.CourseType, error)

	// --- classroom <-> classroom.Classroom ---

	// ClassroomToPO 领域模型 -> 数据模型（写路径）。
	ClassroomToPO(do *classroom.Classroom) *model.Classroom
	// ClassroomToDO 数据模型 -> 领域模型（读路径）。
	ClassroomToDO(po *model.Classroom) (*classroom.Classroom, error)

	// --- student <-> student.Student ---

	// StudentToPO 领域模型 -> 数据模型（写路径）。
	StudentToPO(do *student.Student) *model.Student
	// StudentToDO 数据模型 -> 领域模型（读路径）。
	StudentToDO(po *model.Student) (*student.Student, error)

	// --- teacher <-> teacher.Teacher ---

	// TeacherToPO 领域模型 -> 数据模型（写路径）。
	TeacherToPO(do *teacher.Teacher) *model.Teacher
	// TeacherToDO 数据模型 -> 领域模型（读路径）。
	TeacherToDO(po *model.Teacher) (*teacher.Teacher, error)

	// --- course <-> course.Course ---

	// CourseToPO 领域模型 -> 数据模型（写路径）。
	CourseToPO(do *course.Course) *model.Course
	// CourseToDO 数据模型 -> 领域模型（读路径）。
	CourseToDO(po *model.Course) (*course.Course, error)

	// --- course_slot <-> courseSlot.CourseSlot ---

	// CourseSlotToPO 领域模型 -> 数据模型（写路径）。
	CourseSlotToPO(do *courseSlot.CourseSlot) *model.CourseSlot
	// CourseSlotToDO 数据模型 -> 领域模型（读路径）。
	CourseSlotToDO(po *model.CourseSlot) (*courseSlot.CourseSlot, error)

	// --- course_slot_change <-> courseSlotChange.CourseSlotChange ---

	// CourseSlotChangeToPO 领域模型 -> 数据模型（写路径）。
	CourseSlotChangeToPO(do *courseSlotChange.CourseSlotChange) *model.CourseSlotChange
	// CourseSlotChangeToDO 数据模型 -> 领域模型（读路径）。
	CourseSlotChangeToDO(po *model.CourseSlotChange) (*courseSlotChange.CourseSlotChange, error)

	// --- course_enrollment <-> enrollment.CourseEnrollment ---

	// EnrollmentToPO 领域模型 -> 数据模型（写路径）。
	EnrollmentToPO(do *enrollment.CourseEnrollment) *model.CourseEnrollment
	// EnrollmentToDO 数据模型 -> 领域模型（读路径）。
	EnrollmentToDO(po *model.CourseEnrollment) (*enrollment.CourseEnrollment, error)

	// --- absence_record <-> absence.AbsenceRecord ---

	// AbsenceToPO 领域模型 -> 数据模型（写路径）。
	AbsenceToPO(do *absence.AbsenceRecord) *model.AbsenceRecord
	// AbsenceToDO 数据模型 -> 领域模型（读路径）。
	AbsenceToDO(po *model.AbsenceRecord) (*absence.AbsenceRecord, error)

	// --- student_makeup <-> makeup.StudentMakeup ---

	// MakeupToPO 领域模型 -> 数据模型（写路径）。
	MakeupToPO(do *makeup.StudentMakeup) *model.StudentMakeup
	// MakeupToDO 数据模型 -> 领域模型（读路径）。
	MakeupToDO(po *model.StudentMakeup) (*makeup.StudentMakeup, error)

	// --- qualification <-> qualification.Qualification ---

	// QualificationToPO 领域模型 -> 数据模型（写路径）。
	QualificationToPO(do *qualification.Qualification) *model.Qualification
	// QualificationToDO 数据模型 -> 领域模型（读路径）。
	QualificationToDO(po *model.Qualification) (*qualification.Qualification, error)
}
