// Package memory 是 course_scheduling 读侧的内存实现。
//
// 内存和 MySQL 不一样：这里没有「两个存储」，只有**一份 map**，写进去立刻
// 就要能读到。所以内存的读写分离是**接口实现分离**，不是存储分离 ——
// 完整的那份数据在 adapters/memorystore，两侧共用同一个实例。
//
// 本包只做一件事：把完整存储收窄成**只读面**再交给 implement/ 里的仓库。
//
//	reader  只列读侧真正用到的方法
//	Data    仓库拿到的就是这个，写方法压根不在类型里
//
// 于是读侧调 SaveCourse / DeleteStudent 这类写操作会直接在编译期报错，
// 而不是靠「大家自觉」。这也是本包存在的唯一理由。
package memory

import (
	"cqrs/internal/core/course_scheduling/adapters/memorystore"
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

// reader 是读侧需要的那部分内存存储能力。
//
// 只列 implement/ 里真正调到的方法：读侧只用复数 getter 取整表，
// 连 XxxByID 都没用到，更不用说 Seed* / Save* / Delete* —— 一概不放进来。
type reader interface {
	Courses() []*course.Course
	CourseTypes() []*courseType.CourseType
	Classrooms() []*classroom.Classroom
	Students() []*student.Student
	Teachers() []*teacher.Teacher
	CourseSlots() []*courseSlot.CourseSlot
	CourseSlotChanges() []*courseSlotChange.CourseSlotChange
	Absences() []*absence.AbsenceRecord
	Enrollments() []*enrollment.CourseEnrollment
	Makeups() []*makeup.StudentMakeup
	Qualifications() []*qualification.Qualification
}

// Data 是读侧仓库持有的存储访问面：只有 reader 里的那些方法。
type Data struct {
	reader
}

// NewData 把完整的内存存储收窄成只读面。
//
// 参数是具体类型、字段是窄接口：收窄在赋值这一行由编译器检查，
// 不满足就会在这里报错，也不会给 wire 添额外的 bind。
func NewData(store *memorystore.Data) *Data {
	return &Data{reader: store}
}
