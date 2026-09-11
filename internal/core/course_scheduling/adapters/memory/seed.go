package memory

import (
	"fmt"
	"time"

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
)

// 演示数据固定使用 2026 年秋季学期，保证每次 SeedDemo 结果完全一致
// （分页顺序、冲突判定都可复现）。
const demoYear = 2026

var (
	demoNow        = time.Date(demoYear, time.September, 1, 10, 0, 0, 0, time.Local)
	demoTermStart  = time.Date(demoYear, time.September, 7, 0, 0, 0, 0, time.Local)
	demoTermEnd    = time.Date(demoYear+1, time.January, 15, 0, 0, 0, 0, time.Local)
	demoEnrollFrom = time.Date(demoYear, time.August, 1, 0, 0, 0, 0, time.Local)
	demoEnrollTo   = time.Date(demoYear, time.September, 6, 0, 0, 0, 0, time.Local)
	demoDropBy     = time.Date(demoYear, time.September, 20, 0, 0, 0, 0, time.Local)
)

// demoDate 返回固定年份的当天 0 点。
func demoDate(month time.Month, day int) time.Time {
	return time.Date(demoYear, month, day, 0, 0, 0, 0, time.Local)
}

// demoDateTime 返回固定年份的具体时刻。
func demoDateTime(month time.Month, day, hour, minute int) time.Time {
	return time.Date(demoYear, month, day, hour, minute, 0, 0, time.Local)
}

// SeedDemo 注入一套互相引用、自洽的演示数据，覆盖全部 10 个聚合。
//
// 数据刻意同时包含「时间冲突」与「不冲突」的课程，方便直接观察查询结果差异
// （以下是按当前实现得到的预期结果）：
//
//	AvailableForStudent(101 陈晨)   -> C002, C004        // 101 在学 C001，被 C001 的周一/周三 09:00-11:00 挡掉 C003
//	AvailableForTeacher(1 张伟)     -> C002, C004        // 张伟已排 C001(周一/周三 09:00-11:00)
//	AvailableForTeacher(2 李娜)     -> C001, C003, C004  // 李娜只排了 C002(周一 14:00-16:00)
//	AvailableForTeacher(3 王强)     -> C002              // 王强周一/周二/周三/周五均已被占
//	AvailableForClassroom(R101)     -> C002, C004        // R101 在周一/周三 09:00-11:00 被占用
//	AvailableForClassroom(R102)     -> 空                // R102 的周一/周二/周三/周五都被占
//
// 反复调用会按相同 ID 覆盖，不会产生重复数据。
func (d *Data) SeedDemo() error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"teachers", d.seedDemoTeachers},
		{"students", d.seedDemoStudents},
		{"classrooms", d.seedDemoClassrooms},
		{"courses", d.seedDemoCourses},
		{"course slots", d.seedDemoCourseSlots},
		{"enrollments", d.seedDemoEnrollments},
		{"qualifications", d.seedDemoQualifications},
		{"absences", d.seedDemoAbsences},
		{"makeups", d.seedDemoMakeups},
		{"course slot changes", d.seedDemoCourseSlotChanges},
	}

	for _, step := range steps {
		if err := step.fn(); err != nil {
			return fmt.Errorf("seed demo %s: %w", step.name, err)
		}
	}
	return nil
}

// --- 讲师：1 张伟 / 2 李娜 / 3 王强 ---

func (d *Data) seedDemoTeachers() error {
	t1, err := demoTeacher(1, "张伟", "金牌讲师", "13800000001", "zhangwei@example.com")
	if err != nil {
		return err
	}
	t2, err := demoTeacher(2, "李娜", "特级培训师", "13800000002", "lina@example.com")
	if err != nil {
		return err
	}
	t3, err := demoTeacher(3, "王强", "讲师", "13800000003", "wangqiang@example.com")
	if err != nil {
		return err
	}
	d.SeedTeacher(t1, t2, t3)
	return nil
}

// --- 学员：101 陈晨 / 102 刘洋 / 103 赵敏 ---

func (d *Data) seedDemoStudents() error {
	s101, err := demoStudent(101, "陈晨", student.TypeInternal, "13900000101", "chenchen@example.com")
	if err != nil {
		return err
	}
	s102, err := demoStudent(102, "刘洋", student.TypeExternal, "13900000102", "liuyang@example.com")
	if err != nil {
		return err
	}
	s103, err := demoStudent(103, "赵敏", student.TypeExternal, "13900000103", "zhaomin@example.com")
	if err != nil {
		return err
	}
	d.SeedStudent(s101, s102, s103)
	return nil
}

// --- 教室：R101(30 座) / R102(20 座) ---

func (d *Data) seedDemoClassrooms() error {
	r101, err := demoClassroom("R101", "A座", 3, "301", 30)
	if err != nil {
		return err
	}
	r102, err := demoClassroom("R102", "A座", 3, "302", 20)
	if err != nil {
		return err
	}
	d.SeedClassroom(r101, r102)
	return nil
}

// --- 课程：C001..C004 ---

func (d *Data) seedDemoCourses() error {
	c001, err := demoCourse("C001", 30, 2)
	if err != nil {
		return err
	}
	c002, err := demoCourse("C002", 20, 1)
	if err != nil {
		return err
	}
	c003, err := demoCourse("C003", 15, 0)
	if err != nil {
		return err
	}
	c004, err := demoCourse("C004", 12, 0)
	if err != nil {
		return err
	}
	d.SeedCourse(c001, c002, c003, c004)
	return nil
}

// --- 周排期：刻意制造「周一/周三 09:00-11:00」这段热点时间 ---
//
//	1 C001 周一 09:00-11:00 张伟 R101
//	2 C001 周三 09:00-11:00 张伟 R101
//	3 C002 周一 14:00-16:00 李娜 R102
//	4 C003 周三 09:00-11:00 王强 R102   <- 与 2 同时间：张伟/王强互相撞车
//	5 C003 周五 09:00-11:00 王强 R102
//	6 C004 周二 09:00-11:00 王强 R102
func (d *Data) seedDemoCourseSlots() error {
	slot1, err := demoSlot("550e8400-e29b-41d4-a716-446655440001", "C001", time.Monday, 9, 0, 11, 0, 1, "R101")
	if err != nil {
		return err
	}
	slot2, err := demoSlot("550e8400-e29b-41d4-a716-446655440002", "C001", time.Wednesday, 9, 0, 11, 0, 1, "R101")
	if err != nil {
		return err
	}
	slot3, err := demoSlot("550e8400-e29b-41d4-a716-446655440003", "C002", time.Monday, 14, 0, 16, 0, 2, "R102")
	if err != nil {
		return err
	}
	slot4, err := demoSlot("550e8400-e29b-41d4-a716-446655440004", "C003", time.Wednesday, 9, 0, 11, 0, 3, "R102")
	if err != nil {
		return err
	}
	slot5, err := demoSlot("550e8400-e29b-41d4-a716-446655440005", "C003", time.Friday, 9, 0, 11, 0, 3, "R102")
	if err != nil {
		return err
	}
	slot6, err := demoSlot("550e8400-e29b-41d4-a716-446655440006", "C004", time.Tuesday, 9, 0, 11, 0, 3, "R102")
	if err != nil {
		return err
	}
	d.SeedCourseSlot(slot1, slot2, slot3, slot4, slot5, slot6)
	return nil
}

// --- 注册：103 有一条「已结业」记录，用来验证查询侧不过滤状态 ---
func (d *Data) seedDemoEnrollments() error {
	now := time.Now()
	// 使用 Reconstitute 恢复固定 ID的种子数据
	e1 := enrollment.Reconstitute(1, 101, "C001", enrollment.StatusEnrolled, now, time.Time{}, time.Time{}, now)
	e2 := enrollment.Reconstitute(2, 102, "C002", enrollment.StatusEnrolled, now, time.Time{}, time.Time{}, now)
	e3 := enrollment.Reconstitute(3, 103, "C001", enrollment.StatusEnrolled, now, time.Time{}, time.Time{}, now)
	e4 := enrollment.Reconstitute(4, 103, "C002", enrollment.StatusEnrolled, now, time.Time{}, time.Time{}, now)
	// 已结业：CoursesByStudentID(103) 仍会返回 C002（不过滤状态）
	if err := e4.Complete(demoNow); err != nil {
		return err
	}
	d.SeedEnrollment(e1, e2, e3, e4)
	return nil
}

// --- 资质：其中 1-张伟-C004 被吊销，用来验证查询侧不过滤状态 ---
func (d *Data) seedDemoQualifications() error {
	q1, err := demoQualification(1, 1, "C001")
	if err != nil {
		return err
	}
	q2, err := demoQualification(2, 2, "C002")
	if err != nil {
		return err
	}
	q3, err := demoQualification(3, 3, "C003")
	if err != nil {
		return err
	}
	q4, err := demoQualification(4, 3, "C004")
	if err != nil {
		return err
	}
	q5, err := demoQualification(5, 1, "C004")
	if err != nil {
		return err
	}
	q5.Revoke() // CoursesByTeacherID(1) 仍会返回 C004

	d.SeedQualification(q1, q2, q3, q4, q5)
	return nil
}

// --- 缺勤：1 条待审批事假 + 1 条已生效旷课 ---
func (d *Data) seedDemoAbsences() error {
	now := time.Now()
	// 使用 Reconstitute 恢复固定 ID的种子数据
	a1 := absence.Reconstitute(
		1, 101, "C001", 2, demoDate(time.September, 16), 2,
		absence.TypePersonalLeave, "家中急事",
		absence.StatusPending, 0, "", false, now, now,
	)
	a2 := absence.Reconstitute(
		2, 102, "C002", 3, demoDate(time.September, 14), 2,
		absence.TypeUnexcused, "未请假缺席",
		absence.StatusApproved, 0, "", false, now, now,
	)
	d.SeedAbsence(a1, a2)
	return nil
}

// --- 补课：1 条已核销（走完 审批->完成），1 条待审批 ---
func (d *Data) seedDemoMakeups() error {
	// 使用 Reconstitute 恢复固定 ID的种子数据
	m1 := makeup.Reconstitute(
		1, 101, "C001", 1, 1, demoDate(time.September, 14), 2, demoDate(time.September, 16), 2,
		makeup.StatusApproved, 9, "同意补课", demoNow, demoNow, demoNow,
	)
	if err := m1.CompleteAttendance(demoNow); err != nil {
		return err
	}

	m2 := makeup.Reconstitute(
		2, 102, "C002", 2, 3, demoDate(time.September, 14), 3, demoDate(time.September, 21), 2,
		makeup.StatusPending, 0, "", time.Time{}, demoNow, demoNow,
	)

	d.SeedMakeup(m1, m2)
	return nil
}

// --- 课表变更：1 条已审批的调课 + 1 条待审批的代课 ---
func (d *Data) seedDemoCourseSlotChanges() error {
	original1 := courseSlotChange.NewOriginalPlan(
		1, demoDate(time.September, 14), 1, "R101", "09:00", "11:00",
	)
	target1, err := courseSlotChange.NewTargetPlan(
		demoDateTime(time.September, 15, 14, 0),
		demoDateTime(time.September, 15, 16, 0),
		1, "R102",
	)
	if err != nil {
		return err
	}
	// 使用 Reconstitute 恢复固定 ID的种子数据
	change1 := courseSlotChange.Reconstitute(
		1, "C001", 1, courseSlotChange.TypeReschedule,
		original1, target1, "场地检修，临时调至周二下午",
		courseSlotChange.StatusApproved, 9, "同意调课", demoNow, demoNow,
	)

	original2 := courseSlotChange.NewOriginalPlan(
		2, demoDate(time.September, 16), 1, "R101", "09:00", "11:00",
	)
	target2, err := courseSlotChange.NewTargetPlan(
		demoDateTime(time.September, 16, 9, 0),
		demoDateTime(time.September, 16, 11, 0),
		2, "R101",
	)
	if err != nil {
		return err
	}
	change2 := courseSlotChange.Reconstitute(
		2, "C001", 1, courseSlotChange.TypeSubstitute,
		original2, target2, "讲师出差，由李娜代课",
		courseSlotChange.StatusPending, 0, "", demoNow, demoNow,
	)

	d.SeedCourseSlotChange(change1, change2)
	return nil
}

// --- 构造辅助 ---

func demoTeacher(id int64, name, title, phone, email string) (*teacher.Teacher, error) {
	contact, err := teacher.NewContactInfo(phone, email)
	if err != nil {
		return nil, err
	}
	// 使用 Reconstitute 恢复固定 ID 的种子数据
	return teacher.Reconstitute(id, name, title, contact, teacher.StatusActive, time.Now(), time.Now()), nil
}

func demoStudent(id int64, name string, typ student.StudentType, phone, email string) (*student.Student, error) {
	contact, err := student.NewContactInfo(phone, email)
	if err != nil {
		return nil, err
	}
	// 使用 Reconstitute 恢复固定 ID 的种子数据
	return student.Reconstitute(id, name, typ, contact, student.StatusActive, time.Now(), time.Now()), nil
}

func demoClassroom(id, building string, floor int, room string, capacity int) (*classroom.Classroom, error) {
	location, err := classroom.NewLocation(building, floor, room)
	if err != nil {
		return nil, err
	}
	return classroom.NewClassroom(id, location, capacity)
}

func demoCourse(id string, maxSeats, enrolled int) (*course.Course, error) {
	capacity, err := course.NewCapacity(maxSeats, enrolled)
	if err != nil {
		return nil, err
	}
	period, err := course.NewCoursePeriod(demoTermStart, demoTermEnd, 48, enrolled*4)
	if err != nil {
		return nil, err
	}
	window := course.NewEnrollmentWindow(demoEnrollFrom, demoEnrollTo, demoDropBy)
	return course.NewCourse(id, capacity, window, period), nil
}

func demoSlot(
	id string,
	courseID string,
	weekday time.Weekday,
	fromHour, fromMinute, toHour, toMinute int,
	teacherID int64,
	classroomID string,
) (*courseSlot.CourseSlot, error) {
	from, err := courseSlot.NewDayTime(fromHour, fromMinute)
	if err != nil {
		return nil, err
	}
	to, err := courseSlot.NewDayTime(toHour, toMinute)
	if err != nil {
		return nil, err
	}
	span, err := courseSlot.NewDayTimeRange(from, to)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return courseSlot.Reconstitute(id, courseID, weekday, span, teacherID, classroomID, now, now), nil
}

func demoQualification(id, teacherID int64, courseID string) (*qualification.Qualification, error) {
	// 使用 Reconstitute 恢复固定 ID 的种子数据
	return qualification.Reconstitute(id, teacherID, courseID, demoTermStart, demoTermEnd.AddDate(1, 0, 0), qualification.StatusActive, time.Now()), nil
}
