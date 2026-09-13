package courseSlot

import (
	"testing"
	"time"

	"cqrs/internal/core/course_scheduling/adapters/test_memory"
	commandmemory "cqrs/internal/core/course_scheduling/adapters/test_memory/command"
	memorycmd "cqrs/internal/core/course_scheduling/adapters/test_memory/command/implement"
	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
	"cqrs/internal/core/course_scheduling/domain/service/classroomCapacity"
	"cqrs/internal/core/course_scheduling/domain/service/qualificationCheck"
	"cqrs/internal/core/course_scheduling/domain/service/scheduleConflict"
)

const (
	// ownCourse 课程容量 20 人；otherCourse 只用来占别的时间，不在存储里。
	ownCourse   = "C001"
	otherCourse = "C999"

	// room 装得下 20 人；smallRoom 只有 5 个座，装不下。
	room      = "R101"
	smallRoom = "R102"

	courseSeats  = 20
	roomSeats    = 30
	smallSeats   = 5
	noTeacher    = int64(-1)
	noClassroom  = "-"
	targetPeriod = 2 * time.Hour
)

// fixture 是一份搭好的内存数据 + 处理器。
//
// data 是完整 store（测试靠它塞数据、做断言），handler 里的仓库拿到的
// 则是收窄后的写侧面。
type fixture struct {
	handler   *Handler
	data      *test_memory.Data
	teacher   teacher.Teacher
	course    *course.Course
	room      *classroom.Classroom
	smallRoom *classroom.Classroom
}

// newFixture 造一份数据：
// 1 个课程类型、1 门 20 人的课、1 位讲师、1 间 30 人的教室、1 间 5 人的小教室。
//
// withQualification 决定这位讲师是否持有该课程类型的资质。
func newFixture(t *testing.T, withQualification bool) *fixture {
	t.Helper()

	d, cleanup, err := test_memory.NewData(nil)
	if err != nil {
		t.Fatalf("new data: %v", err)
	}
	t.Cleanup(cleanup)

	contact, err := teacher.NewContactInfo("13800000000", "")
	if err != nil {
		t.Fatalf("new contact: %v", err)
	}
	// ID 固定，所以走 Reconstitute 而不是会自己生成 ID 的 NewTeacher
	teach := teacher.Reconstitute(1, 201, "李娜", "讲师", contact, teacher.StatusActive, time.Now(), time.Now())

	ct, err := courseType.NewCourseType("少儿编程", "")
	if err != nil {
		t.Fatalf("new course type: %v", err)
	}

	capacity, err := course.NewCapacity(courseSeats, 0)
	if err != nil {
		t.Fatalf("new capacity: %v", err)
	}
	now := time.Now()
	period, err := course.NewCoursePeriod(now.AddDate(0, 0, -7), now.AddDate(0, 6, 0), 16, 0)
	if err != nil {
		t.Fatalf("new period: %v", err)
	}
	enrollment := course.NewEnrollmentWindow(
		now.AddDate(0, 0, -1), now.AddDate(0, 0, 7), now.AddDate(0, 0, 14),
	)
	crs := course.Reconstitute(ownCourse, ct.ID(), capacity, enrollment, period)

	big, err := newClassroom(t, room, roomSeats)
	if err != nil {
		t.Fatalf("new classroom: %v", err)
	}
	small, err := newClassroom(t, smallRoom, smallSeats)
	if err != nil {
		t.Fatalf("new classroom: %v", err)
	}

	d.SeedTeacher(teach)
	d.SeedCourseType(ct)
	d.SeedCourse(crs)
	d.SeedClassroom(big, small)

	if withQualification {
		q := qualification.Reconstitute(
			1, teach.ID(), ct.ID(),
			now.AddDate(0, 0, -30),
			qualification.StatusActive, now,
		)
		d.SeedQualification(q)
	}

	data := commandmemory.NewData(d)
	return &fixture{
		handler: NewHandler(
			memorycmd.NewCourseSlotCommand(data),
			scheduleConflict.NewService(
				memorycmd.NewCourseSlotCommand(data),
				memorycmd.NewTeacherCommand(data),
				memorycmd.NewCourseCommand(data),
				memorycmd.NewCourseSlotChangeCommand(data),
			),
			qualificationCheck.NewService(
				memorycmd.NewCourseSlotCommand(data),
				memorycmd.NewTeacherCommand(data),
				memorycmd.NewQualificationCommand(data),
				memorycmd.NewCourseTypeCommand(data),
				memorycmd.NewCourseCommand(data),
				memorycmd.NewEnrollmentCommand(data),
				memorycmd.NewAbsenceCommand(data),
			),
			classroomCapacity.NewService(
				memorycmd.NewClassroomCommand(data),
				memorycmd.NewCourseSlotCommand(data),
				memorycmd.NewCourseCommand(data),
				memorycmd.NewMakeupCommand(data),
			),
			data,
		),
		data:      d,
		teacher:   *teach,
		course:    crs,
		room:      big,
		smallRoom: small,
	}
}

// newClassroom 造一间教室，location 用 ID 拼出来；ID 固定所以走 Reconstitute。
func newClassroom(t *testing.T, id string, seats int) (*classroom.Classroom, error) {
	t.Helper()

	loc, err := classroom.NewLocation("A", 1, id)
	if err != nil {
		t.Fatalf("new location: %v", err)
	}
	return classroom.Reconstitute(id, loc, seats, 0, classroom.StatusAvailable), nil
}

// newSlot 造一条待排的槽位（还没排教室）并塞进存储。
func (f *fixture) newSlot(t *testing.T, courseID string, teacherID int64, weekday time.Weekday, fromHour, toHour int) *courseSlot.CourseSlot {
	t.Helper()
	return f.newSlotIn(t, courseID, teacherID, noClassroom, weekday, fromHour, toHour)
}

// newSlotIn 造一条待排的槽位（指定教室）并塞进存储。
func (f *fixture) newSlotIn(t *testing.T, courseID string, teacherID int64, classroomID string, weekday time.Weekday, fromHour, toHour int) *courseSlot.CourseSlot {
	t.Helper()

	from, err := courseSlot.NewDayTime(fromHour, 0)
	if err != nil {
		t.Fatalf("new day time: %v", err)
	}
	to, err := courseSlot.NewDayTime(toHour, 0)
	if err != nil {
		t.Fatalf("new day time: %v", err)
	}
	span, err := courseSlot.NewDayTimeRange(from, to)
	if err != nil {
		t.Fatalf("new day time range: %v", err)
	}

	cs, err := courseSlot.NewCourseSlot(courseID, weekday, span, teacherID, classroomID)
	if err != nil {
		t.Fatalf("new course slot: %v", err)
	}
	f.data.SeedCourseSlot(cs)

	return cs
}

// slotTeacherID 读回槽位当前的讲师 ID。
func (f *fixture) slotTeacherID(t *testing.T, id string) int64 {
	t.Helper()

	cs, ok := f.data.CourseSlotByID(id)
	if !ok {
		t.Fatalf("slot %s not found", id)
	}
	return cs.TeacherID()
}

// slotCourseID 读回槽位当前的课程 ID。
func (f *fixture) slotCourseID(t *testing.T, id string) string {
	t.Helper()

	cs, ok := f.data.CourseSlotByID(id)
	if !ok {
		t.Fatalf("slot %s not found", id)
	}
	return cs.CourseID()
}

// slotClassroomID 读回槽位当前的教室 ID。
func (f *fixture) slotClassroomID(t *testing.T, id string) string {
	t.Helper()

	cs, ok := f.data.CourseSlotByID(id)
	if !ok {
		t.Fatalf("slot %s not found", id)
	}
	return cs.ClassroomID()
}
