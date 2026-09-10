package memory

import (
	"sort"
	"sync"

	"cqrs/internal/conf"
	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"

	"github.com/go-kratos/kratos/v3/log"
)

// Data 是 course_scheduling 上下文的内存存储，command / query 适配器共享同一份。
//
// 目前没有接持久化：数据只存活在进程内，重启即清空。
// 查询适配器只读；数据通过 Seed* 注入（测试、本地演示），
// 真实写路径等 app/command 落地后再补。
type Data struct {
	mu sync.RWMutex

	courses           map[string]*course.Course
	classrooms        map[string]*classroom.Classroom
	students          map[int64]*student.Student
	teachers          map[int64]*teacher.Teacher
	courseSlots       map[int64]*courseSlot.CourseSlot
	courseSlotChanges map[int64]*courseSlotChange.CourseSlotChange
	absences          map[int64]*absence.AbsenceRecord
	enrollments       map[int64]*enrollment.CourseEnrollment
	makeups           map[int64]*makeup.StudentMakeup
	qualifications    map[int64]*qualification.Qualification
}

// NewData 创建一个空的内存存储。
func NewData(_ *conf.Data) (*Data, func(), error) {
	d := &Data{
		courses:           make(map[string]*course.Course),
		classrooms:        make(map[string]*classroom.Classroom),
		students:          make(map[int64]*student.Student),
		teachers:          make(map[int64]*teacher.Teacher),
		courseSlots:       make(map[int64]*courseSlot.CourseSlot),
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
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, item := range items {
		if item != nil {
			d.courses[item.ID()] = item
		}
	}
}

// SeedClassroom 写入或覆盖教室。
func (d *Data) SeedClassroom(items ...*classroom.Classroom) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, item := range items {
		if item != nil {
			d.classrooms[item.ID()] = item
		}
	}
}

// SeedStudent 写入或覆盖学员。
func (d *Data) SeedStudent(items ...*student.Student) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, item := range items {
		if item != nil {
			d.students[item.ID()] = item
		}
	}
}

// SeedTeacher 写入或覆盖讲师。
func (d *Data) SeedTeacher(items ...*teacher.Teacher) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, item := range items {
		if item != nil {
			d.teachers[item.ID()] = item
		}
	}
}

// SeedCourseSlot 写入或覆盖课表槽位。
func (d *Data) SeedCourseSlot(items ...*courseSlot.CourseSlot) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, item := range items {
		if item != nil {
			d.courseSlots[item.ID()] = item
		}
	}
}

// SeedCourseSlotChange 写入或覆盖课表变更单。
func (d *Data) SeedCourseSlotChange(items ...*courseSlotChange.CourseSlotChange) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, item := range items {
		if item != nil {
			d.courseSlotChanges[item.ID()] = item
		}
	}
}

// SeedAbsence 写入或覆盖缺勤记录。
func (d *Data) SeedAbsence(items ...*absence.AbsenceRecord) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, item := range items {
		if item != nil {
			d.absences[item.ID()] = item
		}
	}
}

// SeedEnrollment 写入或覆盖课程注册记录。
func (d *Data) SeedEnrollment(items ...*enrollment.CourseEnrollment) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, item := range items {
		if item != nil {
			d.enrollments[item.ID()] = item
		}
	}
}

// SeedMakeup 写入或覆盖补课申请。
func (d *Data) SeedMakeup(items ...*makeup.StudentMakeup) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, item := range items {
		if item != nil {
			d.makeups[item.ID()] = item
		}
	}
}

// SeedQualification 写入或覆盖授课资质。
func (d *Data) SeedQualification(items ...*qualification.Qualification) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, item := range items {
		if item != nil {
			d.qualifications[item.ID()] = item
		}
	}
}

// --- 读取：返回按 ID 升序的快照，调用方拿到的是新切片 ---

// Courses 返回全部课程快照（按课程 ID 升序）。
func (d *Data) Courses() []*course.Course {
	d.mu.RLock()
	defer d.mu.RUnlock()
	out := make([]*course.Course, 0, len(d.courses))
	for _, item := range d.courses {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

// Classrooms 返回全部教室快照（按教室 ID 升序）。
func (d *Data) Classrooms() []*classroom.Classroom {
	d.mu.RLock()
	defer d.mu.RUnlock()
	out := make([]*classroom.Classroom, 0, len(d.classrooms))
	for _, item := range d.classrooms {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

// Students 返回全部学员快照（按学员 ID 升序）。
func (d *Data) Students() []*student.Student {
	d.mu.RLock()
	defer d.mu.RUnlock()
	out := make([]*student.Student, 0, len(d.students))
	for _, item := range d.students {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

// Teachers 返回全部讲师快照（按讲师 ID 升序）。
func (d *Data) Teachers() []*teacher.Teacher {
	d.mu.RLock()
	defer d.mu.RUnlock()
	out := make([]*teacher.Teacher, 0, len(d.teachers))
	for _, item := range d.teachers {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

// CourseSlots 返回全部课表槽位快照（按槽位 ID 升序）。
func (d *Data) CourseSlots() []*courseSlot.CourseSlot {
	d.mu.RLock()
	defer d.mu.RUnlock()
	out := make([]*courseSlot.CourseSlot, 0, len(d.courseSlots))
	for _, item := range d.courseSlots {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

// CourseSlotChanges 返回全部课表变更快照（按变更单 ID 升序）。
func (d *Data) CourseSlotChanges() []*courseSlotChange.CourseSlotChange {
	d.mu.RLock()
	defer d.mu.RUnlock()
	out := make([]*courseSlotChange.CourseSlotChange, 0, len(d.courseSlotChanges))
	for _, item := range d.courseSlotChanges {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

// Absences 返回全部缺勤记录快照（按记录 ID 升序）。
func (d *Data) Absences() []*absence.AbsenceRecord {
	d.mu.RLock()
	defer d.mu.RUnlock()
	out := make([]*absence.AbsenceRecord, 0, len(d.absences))
	for _, item := range d.absences {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

// Enrollments 返回全部课程注册快照（按注册 ID 升序）。
func (d *Data) Enrollments() []*enrollment.CourseEnrollment {
	d.mu.RLock()
	defer d.mu.RUnlock()
	out := make([]*enrollment.CourseEnrollment, 0, len(d.enrollments))
	for _, item := range d.enrollments {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

// Makeups 返回全部补课申请快照（按申请 ID 升序）。
func (d *Data) Makeups() []*makeup.StudentMakeup {
	d.mu.RLock()
	defer d.mu.RUnlock()
	out := make([]*makeup.StudentMakeup, 0, len(d.makeups))
	for _, item := range d.makeups {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}

// Qualifications 返回全部授课资质快照（按资质 ID 升序）。
func (d *Data) Qualifications() []*qualification.Qualification {
	d.mu.RLock()
	defer d.mu.RUnlock()
	out := make([]*qualification.Qualification, 0, len(d.qualifications))
	for _, item := range d.qualifications {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID() < out[j].ID() })
	return out
}
