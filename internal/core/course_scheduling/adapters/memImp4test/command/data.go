// Package memory 是 course_scheduling 写侧的内存实现。
//
// 内存和 MySQL 不一样：这里没有「两个存储」，只有**一份 map**，写进去立刻
// 就要能读到。所以内存的读写分离是**接口实现分离**，不是存储分离 ——
// 完整的那份数据在 adapters/memImp4test，两侧共用同一个实例。
//
// 本包只做一件事：把完整存储收窄成写侧需要的那一面再交给 implement/ 里的仓库。
//
//	writer  写方法（Save* / Delete*）
//	reader  读方法；写侧**必须**能读（「先读排期判冲突、再写回」就是这个形状）
//	Data    同时嵌入两者，仓库拿到的就是它
//
// 读侧只有 reader，没有 writer —— 写方法在那里编译期就不可见。
package command

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/memImp4test"
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

// reader 是写侧读数据需要的那部分能力。
//
// 写操作几乎都是「先查再改」：判冲突要读排期、加人要读课程与名额，
// 所以写侧比读侧多用到一批按 ID 取单个的方法。
type reader interface {
	Courses() []*course.Course
	CourseTypes() []*courseType.CourseType
	Classrooms() []*classroom.Classroom
	Teachers() []*teacher.Teacher
	CourseSlots() []*courseSlot.CourseSlot
	CourseSlotChanges() []*courseSlotChange.CourseSlotChange
	Absences() []*absence.AbsenceRecord
	Enrollments() []*enrollment.CourseEnrollment
	Qualifications() []*qualification.Qualification
	Makeups() []*makeup.StudentMakeup

	CourseByID(id string) (*course.Course, bool)
	ClassroomByID(id string) (*classroom.Classroom, bool)
	StudentByID(id int64) (*student.Student, bool)
	TeacherByID(id int64) (*teacher.Teacher, bool)
	CourseSlotByID(id string) (*courseSlot.CourseSlot, bool)
}

// writer 是写侧改数据需要的那部分能力。
//
// 只有 Save* / Delete*：聚合根先在领域模型里改好，再整条覆盖回去
// （内存里存的就是领域对象本身，没有 PO 这一层）。
// Seed* 不在里面 —— 那是测试和演示用的注入手段，属于 memory 的事。
type writer interface {
	SaveCourse(c *course.Course)
	DeleteCourse(id string) bool
	SaveClassroom(c *classroom.Classroom)
	DeleteClassroom(id string) bool
	SaveStudent(s *student.Student)
	DeleteStudent(id int64) bool
	SaveTeacher(t *teacher.Teacher)
	DeleteTeacher(id int64) bool
	SaveCourseSlot(cs *courseSlot.CourseSlot)
	DeleteCourseSlot(id string) bool
	SaveCourseSlotChange(csc *courseSlotChange.CourseSlotChange)
	DeleteCourseSlotChange(id int64) bool
	SaveAbsence(a *absence.AbsenceRecord)
	DeleteAbsence(id int64) bool
	SaveEnrollment(e *enrollment.CourseEnrollment)
	DeleteEnrollment(id int64) bool
	SaveMakeup(m *makeup.StudentMakeup)
	DeleteMakeup(id int64) bool
	SaveQualification(q *qualification.Qualification)
	DeleteQualification(id int64) bool
}

// Data 是写侧仓库持有的存储访问面：能读也能写。
type Data struct {
	reader
	writer
}

// NewData 把完整的内存存储收窄成写侧要的那一面。
//
// 参数是具体类型、字段是窄接口：收窄在赋值这一行由编译器检查，
// 不满足就会在这里报错，也不会给 wire 添额外的 bind。
func NewData(store *memImp4test.Data) *Data {
	return &Data{reader: store, writer: store}
}

// Begin 满足 app/command.Transaction。
//
// 内存里没有真事务：就一份 map，不存在「写到一半失败要回滚」这回事，所以 ctx
// 原样返回。它的意义是让需要事务的命令处理器在内存 profile 下也装得起来，
// 且「判定 + 写回」那段代码在两种实现下写法一致。
func (d *Data) Begin(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

// End 满足 app/command.Transaction；内存里没什么可收尾的，把 err 透出去。
func (d *Data) End(_ context.Context, err error) error {
	return err
}
