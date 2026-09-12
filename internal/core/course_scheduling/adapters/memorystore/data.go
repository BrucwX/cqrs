package memorystore

import (
	"sort"
	"sync"

	"cqrs/internal/conf"
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

	"github.com/go-kratos/kratos/v3/log"
)

// Data 是 course_scheduling 上下文的内存存储，command / query 适配器共享同一份。
//
// 类比数据库：每个 map 是一张"表"，各带一把表锁（<xxx>Mu）；
// 读写只锁自己那张表，不同表之间互不阻塞。
//
// 目前没有接持久化：数据只存活在进程内，重启即清空。
// 查询适配器只读；数据通过 Seed* 注入（测试、本地演示），
// 真实写路径等 app/command 落地后再补。
type Data struct {
	coursesMu sync.RWMutex
	courses   map[string]*course.Course

	courseTypesMu sync.RWMutex
	courseTypes   map[string]*courseType.CourseType

	classroomsMu sync.RWMutex
	classrooms   map[string]*classroom.Classroom

	studentsMu sync.RWMutex
	students   map[int64]*student.Student

	teachersMu sync.RWMutex
	teachers   map[int64]*teacher.Teacher

	courseSlotsMu sync.RWMutex
	courseSlots   map[string]*courseSlot.CourseSlot

	courseSlotChangesMu sync.RWMutex
	courseSlotChanges   map[int64]*courseSlotChange.CourseSlotChange

	absencesMu sync.RWMutex
	absences   map[int64]*absence.AbsenceRecord

	enrollmentsMu sync.RWMutex
	enrollments   map[int64]*enrollment.CourseEnrollment

	makeupsMu sync.RWMutex
	makeups   map[int64]*makeup.StudentMakeup

	qualificationsMu sync.RWMutex
	qualifications   map[int64]*qualification.Qualification
}

// NewData 创建一个空的内存存储。
func NewData(_ *conf.Data) (*Data, func(), error) {
	d := &Data{
		courses:           make(map[string]*course.Course),
		courseTypes:       make(map[string]*courseType.CourseType),
		classrooms:        make(map[string]*classroom.Classroom),
		students:          make(map[int64]*student.Student),
		teachers:          make(map[int64]*teacher.Teacher),
		courseSlots:       make(map[string]*courseSlot.CourseSlot),
		courseSlotChanges: make(map[int64]*courseSlotChange.CourseSlotChange),
		absences:          make(map[int64]*absence.AbsenceRecord),
		enrollments:       make(map[int64]*enrollment.CourseEnrollment),
		makeups:           make(map[int64]*makeup.StudentMakeup),
		qualifications:    make(map[int64]*qualification.Qualification),
	}

	cleanup := func() {
		log.Info("closing course_scheduling memory adapter")
	}
	return d, cleanup, nil
}

// --- 写入（Seed）：仅用于测试 / 本地演示注入数据 ---

// SeedCourse 写入或覆盖课程。
func (d *Data) SeedCourse(items ...*course.Course) {
	d.coursesMu.Lock()
	defer d.coursesMu.Unlock()
	for _, item := range items {
		if item != nil {
			d.courses[item.ID()] = item
		}
	}
}

// SaveCourse 写入或覆盖单门课程（写侧适配器使用）。
func (d *Data) SaveCourse(c *course.Course) {
	d.SeedCourse(c)
}

// DeleteCourse 按 ID 删除课程，返回是否真的删掉了。
func (d *Data) DeleteCourse(id string) bool {
	d.coursesMu.Lock()
	defer d.coursesMu.Unlock()
	if _, ok := d.courses[id]; !ok {
		return false
	}
	delete(d.courses, id)
	return true
}

// SeedCourseType 写入或覆盖课程类型。
func (d *Data) SeedCourseType(items ...*courseType.CourseType) {
	d.courseTypesMu.Lock()
	defer d.courseTypesMu.Unlock()
	for _, item := range items {
		if item != nil {
			d.courseTypes[item.ID()] = item
		}
	}
}

// SeedClassroom 写入或覆盖教室。
func (d *Data) SeedClassroom(items ...*classroom.Classroom) {
	d.classroomsMu.Lock()
	defer d.classroomsMu.Unlock()
	for _, item := range items {
		if item != nil {
			d.classrooms[item.ID()] = item
		}
	}
}

// SaveClassroom 写入或覆盖单间教室（写侧适配器使用）。
func (d *Data) SaveClassroom(c *classroom.Classroom) {
	d.SeedClassroom(c)
}

// DeleteClassroom 按 ID 删除教室，返回是否真的删掉了。
func (d *Data) DeleteClassroom(id string) bool {
	d.classroomsMu.Lock()
	defer d.classroomsMu.Unlock()
	if _, ok := d.classrooms[id]; !ok {
		return false
	}
	delete(d.classrooms, id)
	return true
}

// SeedStudent 写入或覆盖学员。
func (d *Data) SeedStudent(items ...*student.Student) {
	d.studentsMu.Lock()
	defer d.studentsMu.Unlock()
	for _, item := range items {
		if item != nil {
			d.students[item.ID()] = item
		}
	}
}

// SaveStudent 写入或覆盖单个学员（写侧适配器使用）。
func (d *Data) SaveStudent(s *student.Student) {
	d.SeedStudent(s)
}

// DeleteStudent 按 ID 删除学员，返回是否真的删掉了。
func (d *Data) DeleteStudent(id int64) bool {
	d.studentsMu.Lock()
	defer d.studentsMu.Unlock()
	if _, ok := d.students[id]; !ok {
		return false
	}
	delete(d.students, id)
	return true
}

// SeedTeacher 写入或覆盖讲师。
func (d *Data) SeedTeacher(items ...*teacher.Teacher) {
	d.teachersMu.Lock()
	defer d.teachersMu.Unlock()
	for _, item := range items {
		if item != nil {
			d.teachers[item.ID()] = item
		}
	}
}

// SaveTeacher 写入或覆盖单个讲师（写侧适配器使用）。
func (d *Data) SaveTeacher(t *teacher.Teacher) {
	d.SeedTeacher(t)
}

// DeleteTeacher 按 ID 删除讲师，返回是否真的删掉了。
func (d *Data) DeleteTeacher(id int64) bool {
	d.teachersMu.Lock()
	defer d.teachersMu.Unlock()
	if _, ok := d.teachers[id]; !ok {
		return false
	}
	delete(d.teachers, id)
	return true
}

// SeedCourseSlot 写入或覆盖课表槽位。
func (d *Data) SeedCourseSlot(items ...*courseSlot.CourseSlot) {
	d.courseSlotsMu.Lock()
	defer d.courseSlotsMu.Unlock()
	for _, item := range items {
		if item != nil {
			d.courseSlots[item.ID()] = item
		}
	}
}

// SaveCourseSlot 写入或覆盖单个课表槽位（写侧适配器使用）。
func (d *Data) SaveCourseSlot(cs *courseSlot.CourseSlot) {
	d.SeedCourseSlot(cs)
}

// DeleteCourseSlot 按 ID 删除课表槽位，返回是否真的删掉了。
func (d *Data) DeleteCourseSlot(id string) bool {
	d.courseSlotsMu.Lock()
	defer d.courseSlotsMu.Unlock()
	if _, ok := d.courseSlots[id]; !ok {
		return false
	}
	delete(d.courseSlots, id)
	return true
}

// CourseSlotByID 按 ID 取课表槽位。
func (d *Data) CourseSlotByID(id string) (*courseSlot.CourseSlot, bool) {
	d.courseSlotsMu.RLock()
	defer d.courseSlotsMu.RUnlock()
	item, ok := d.courseSlots[id]
	return item, ok
}

// SeedCourseSlotChange 写入或覆盖课表变更单。
func (d *Data) SeedCourseSlotChange(items ...*courseSlotChange.CourseSlotChange) {
	d.courseSlotChangesMu.Lock()
	defer d.courseSlotChangesMu.Unlock()
	for _, item := range items {
		if item != nil {
			d.courseSlotChanges[item.ID()] = item
		}
	}
}

// SaveCourseSlotChange 写入或覆盖单张课表变更单（写侧适配器使用）。
func (d *Data) SaveCourseSlotChange(csc *courseSlotChange.CourseSlotChange) {
	d.SeedCourseSlotChange(csc)
}

// DeleteCourseSlotChange 按 ID 删除课表变更单，返回是否真的删掉了。
func (d *Data) DeleteCourseSlotChange(id int64) bool {
	d.courseSlotChangesMu.Lock()
	defer d.courseSlotChangesMu.Unlock()
	if _, ok := d.courseSlotChanges[id]; !ok {
		return false
	}
	delete(d.courseSlotChanges, id)
	return true
}

// SeedAbsence 写入或覆盖缺勤记录。
func (d *Data) SeedAbsence(items ...*absence.AbsenceRecord) {
	d.absencesMu.Lock()
	defer d.absencesMu.Unlock()
	for _, item := range items {
		if item != nil {
			d.absences[item.ID()] = item
		}
	}
}

// SaveAbsence 写入或覆盖单条缺勤记录（写侧适配器使用）。
func (d *Data) SaveAbsence(a *absence.AbsenceRecord) {
	d.SeedAbsence(a)
}

// DeleteAbsence 按 ID 删除缺勤记录，返回是否真的删掉了。
func (d *Data) DeleteAbsence(id int64) bool {
	d.absencesMu.Lock()
	defer d.absencesMu.Unlock()
	if _, ok := d.absences[id]; !ok {
		return false
	}
	delete(d.absences, id)
	return true
}

// SeedEnrollment 写入或覆盖课程注册记录。
func (d *Data) SeedEnrollment(items ...*enrollment.CourseEnrollment) {
	d.enrollmentsMu.Lock()
	defer d.enrollmentsMu.Unlock()
	for _, item := range items {
		if item != nil {
			d.enrollments[item.ID()] = item
		}
	}
}

// SaveEnrollment 写入或覆盖单条课程注册（写侧适配器使用）。
func (d *Data) SaveEnrollment(e *enrollment.CourseEnrollment) {
	d.SeedEnrollment(e)
}

// DeleteEnrollment 按 ID 删除课程注册，返回是否真的删掉了。
func (d *Data) DeleteEnrollment(id int64) bool {
	d.enrollmentsMu.Lock()
	defer d.enrollmentsMu.Unlock()
	if _, ok := d.enrollments[id]; !ok {
		return false
	}
	delete(d.enrollments, id)
	return true
}

// SeedMakeup 写入或覆盖补课申请。
func (d *Data) SeedMakeup(items ...*makeup.StudentMakeup) {
	d.makeupsMu.Lock()
	defer d.makeupsMu.Unlock()
	for _, item := range items {
		if item != nil {
			d.makeups[item.ID()] = item
		}
	}
}

// SaveMakeup 写入或覆盖单条补课预约（写侧适配器使用）。
func (d *Data) SaveMakeup(m *makeup.StudentMakeup) {
	d.SeedMakeup(m)
}

// DeleteMakeup 按 ID 删除补课预约，返回是否真的删掉了。
func (d *Data) DeleteMakeup(id int64) bool {
	d.makeupsMu.Lock()
	defer d.makeupsMu.Unlock()
	if _, ok := d.makeups[id]; !ok {
		return false
	}
	delete(d.makeups, id)
	return true
}

// SeedQualification 写入或覆盖授课资质。
func (d *Data) SeedQualification(items ...*qualification.Qualification) {
	d.qualificationsMu.Lock()
	defer d.qualificationsMu.Unlock()
	for _, item := range items {
		if item != nil {
			d.qualifications[item.ID()] = item
		}
	}
}

// SaveQualification 写入或覆盖单条授课资质（写侧适配器使用）。
func (d *Data) SaveQualification(q *qualification.Qualification) {
	d.SeedQualification(q)
}

// DeleteQualification 按 ID 删除授课资质，返回是否真的删掉了。
func (d *Data) DeleteQualification(id int64) bool {
	d.qualificationsMu.Lock()
	defer d.qualificationsMu.Unlock()
	if _, ok := d.qualifications[id]; !ok {
		return false
	}
	delete(d.qualifications, id)
	return true
}

// --- 读取：返回按 ID 升序的快照，调用方拿到的是新切片 ---

// Courses 返回全部课程快照（按课程 ID 升序）。
func (d *Data) Courses() []*course.Course {
	d.coursesMu.RLock()
	defer d.coursesMu.RUnlock()
	out := make([]*course.Course, 0, len(d.courses))
	for _, item := range d.courses {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

// CourseTypes 返回全部课程类型快照（按类型 ID 升序）。
func (d *Data) CourseTypes() []*courseType.CourseType {
	d.courseTypesMu.RLock()
	defer d.courseTypesMu.RUnlock()
	out := make([]*courseType.CourseType, 0, len(d.courseTypes))
	for _, item := range d.courseTypes {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

// Classrooms 返回全部教室快照（按教室 ID 升序）。
func (d *Data) Classrooms() []*classroom.Classroom {
	d.classroomsMu.RLock()
	defer d.classroomsMu.RUnlock()
	out := make([]*classroom.Classroom, 0, len(d.classrooms))
	for _, item := range d.classrooms {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

// ClassroomByID 按 ID 取教室。
func (d *Data) ClassroomByID(id string) (*classroom.Classroom, bool) {
	d.classroomsMu.RLock()
	defer d.classroomsMu.RUnlock()
	item, ok := d.classrooms[id]
	return item, ok
}

// CourseByID 按 ID 取课程。
func (d *Data) CourseByID(id string) (*course.Course, bool) {
	d.coursesMu.RLock()
	defer d.coursesMu.RUnlock()
	item, ok := d.courses[id]
	return item, ok
}

// TeacherByID 按 ID 取讲师。
func (d *Data) TeacherByID(id int64) (*teacher.Teacher, bool) {
	d.teachersMu.RLock()
	defer d.teachersMu.RUnlock()
	item, ok := d.teachers[id]
	return item, ok
}

// StudentByID 按 ID 取学员。
func (d *Data) StudentByID(id int64) (*student.Student, bool) {
	d.studentsMu.RLock()
	defer d.studentsMu.RUnlock()
	item, ok := d.students[id]
	return item, ok
}

// QualificationByID 按 ID 取授课资质。
func (d *Data) QualificationByID(id int64) (*qualification.Qualification, bool) {
	d.qualificationsMu.RLock()
	defer d.qualificationsMu.RUnlock()
	item, ok := d.qualifications[id]
	return item, ok
}

// Students 返回全部学员快照（按学员 ID 升序）。
func (d *Data) Students() []*student.Student {
	d.studentsMu.RLock()
	defer d.studentsMu.RUnlock()
	out := make([]*student.Student, 0, len(d.students))
	for _, item := range d.students {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

// Teachers 返回全部讲师快照（按讲师 ID 升序）。
func (d *Data) Teachers() []*teacher.Teacher {
	d.teachersMu.RLock()
	defer d.teachersMu.RUnlock()
	out := make([]*teacher.Teacher, 0, len(d.teachers))
	for _, item := range d.teachers {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

// CourseSlots 返回全部课表槽位快照（按槽位 ID 升序）。
func (d *Data) CourseSlots() []*courseSlot.CourseSlot {
	d.courseSlotsMu.RLock()
	defer d.courseSlotsMu.RUnlock()
	out := make([]*courseSlot.CourseSlot, 0, len(d.courseSlots))
	for _, item := range d.courseSlots {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

// CourseSlotChanges 返回全部课表变更快照（按变更单 ID 升序）。
func (d *Data) CourseSlotChanges() []*courseSlotChange.CourseSlotChange {
	d.courseSlotChangesMu.RLock()
	defer d.courseSlotChangesMu.RUnlock()
	out := make([]*courseSlotChange.CourseSlotChange, 0, len(d.courseSlotChanges))
	for _, item := range d.courseSlotChanges {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

// Absences 返回全部缺勤记录快照（按记录 ID 升序）。
func (d *Data) Absences() []*absence.AbsenceRecord {
	d.absencesMu.RLock()
	defer d.absencesMu.RUnlock()
	out := make([]*absence.AbsenceRecord, 0, len(d.absences))
	for _, item := range d.absences {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

// Enrollments 返回全部课程注册快照（按注册 ID 升序）。
func (d *Data) Enrollments() []*enrollment.CourseEnrollment {
	d.enrollmentsMu.RLock()
	defer d.enrollmentsMu.RUnlock()
	out := make([]*enrollment.CourseEnrollment, 0, len(d.enrollments))
	for _, item := range d.enrollments {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

// Makeups 返回全部补课申请快照（按申请 ID 升序）。
func (d *Data) Makeups() []*makeup.StudentMakeup {
	d.makeupsMu.RLock()
	defer d.makeupsMu.RUnlock()
	out := make([]*makeup.StudentMakeup, 0, len(d.makeups))
	for _, item := range d.makeups {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

// Qualifications 返回全部授课资质快照（按资质 ID 升序）。
func (d *Data) Qualifications() []*qualification.Qualification {
	d.qualificationsMu.RLock()
	defer d.qualificationsMu.RUnlock()
	out := make([]*qualification.Qualification, 0, len(d.qualifications))
	for _, item := range d.qualifications {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}
