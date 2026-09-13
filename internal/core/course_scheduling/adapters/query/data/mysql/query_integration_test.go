package mysql

import (
	"context"
	"os"
	"slices"
	"testing"

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
)

// 集成测试：对着真实的 course_scheduling 库跑一遍查询。
//
// 库由 script/mysql/course_scheduling.sql + course_scheduling_seed.sql 建好并灌入
// 演示数据，下面的期望值就是那份种子数据（teacher 1-3、student 101-103、
// classroom R101/R102、courseType ct-0001/0002、course C001-C004、slot …440001-440006）。
//
// 容器里 MySQL 的 root 只认 unix socket，所以默认 DSN 用 socket 形式；
// 换环境用 COURSE_SCHEDULING_DSN 覆盖。连不上就跳过，不影响纯单元测试。
const defaultIntegrationDSN = "root@unix(/var/run/mysqld/mysqld.sock)/course_scheduling"

// integrationDSN 返回集成测试要连的库，默认是本机 dev container 的 socket 形式。
func integrationDSN() string {
	if dsn := os.Getenv("COURSE_SCHEDULING_DSN"); dsn != "" {
		return dsn
	}
	return defaultIntegrationDSN
}

// newIntegrationData 连上读库；连不上则跳过整条用例。
func newIntegrationData(t *testing.T) *Data {
	t.Helper()

	data, cleanup, err := NewData(&conf.Data{
		Database: &conf.Data_Database{Driver: "mysql", Source: integrationDSN()},
	})
	if err != nil {
		t.Skipf("跳过集成测试：连不上 MySQL（%v）", err)
	}
	t.Cleanup(cleanup)
	return data
}

// TestReadSource 确认读侧用的确实是 read_source：它指到哪就读哪。
func TestReadSource(t *testing.T) {
	// 留空 -> 回落到 source
	data, cleanup, err := NewData(&conf.Data{Database: &conf.Data_Database{
		Driver: "mysql",
		Source: integrationDSN(),
	}})
	if err != nil {
		t.Fatalf("read_source 留空时应回落到 source: %v", err)
	}
	if data == nil {
		t.Fatal("NewData 返回了 nil Data")
	}
	cleanup()

	// 指向不存在的库 -> 能证明用的确实是 read_source，而不是 source
	if _, _, err := NewData(&conf.Data{Database: &conf.Data_Database{
		Driver:     "mysql",
		Source:     integrationDSN(),
		ReadSource: "root@unix(/var/run/mysqld/mysqld.sock)/course_scheduling_missing",
	}}); err == nil {
		t.Error("read_source 指向不存在的库时，NewData 应该报错")
	}
}

// --- 断言小工具 ---

func strIDs[E any](items []E, id func(E) string) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, id(item))
	}
	return out
}

func intIDs[E any](items []E, id func(E) int64) []int64 {
	out := make([]int64, 0, len(items))
	for _, item := range items {
		out = append(out, id(item))
	}
	return out
}

func courseIDs(items []*course.Course) []string {
	return strIDs(items, func(c *course.Course) string { return c.ID() })
}

func assertStrings(t *testing.T, name string, got, want []string) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}

func assertInts(t *testing.T, name string, got, want []int64) {
	t.Helper()
	if !slices.Equal(got, want) {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}

// --- 分页 ---

func TestIntegrationPage(t *testing.T) {
	data := newIntegrationData(t)
	ctx := context.Background()

	t.Run("course 分页与越界", func(t *testing.T) {
		q := data

		got, err := q.PageCourses(ctx, 1, 2)
		if err != nil {
			t.Fatal(err)
		}
		assertStrings(t, "page 1", courseIDs(got), []string{"C001", "C002"})

		got, err = q.PageCourses(ctx, 2, 2)
		if err != nil {
			t.Fatal(err)
		}
		assertStrings(t, "page 2", courseIDs(got), []string{"C003", "C004"})

		got, err = q.PageCourses(ctx, 9, 2)
		if err != nil {
			t.Fatal(err)
		}
		if got == nil {
			t.Error("越界页必须返回非 nil 空切片")
		}
		assertStrings(t, "page 9", courseIDs(got), []string{})
	})

	t.Run("classroom", func(t *testing.T) {
		got, err := data.PageClassrooms(ctx, 1, 10)
		if err != nil {
			t.Fatal(err)
		}
		assertStrings(t, "classroom", strIDs(got, func(c *classroom.Classroom) string { return c.ID() }),
			[]string{"R101", "R102"})
	})

	t.Run("student", func(t *testing.T) {
		got, err := data.PageStudents(ctx, 1, 2)
		if err != nil {
			t.Fatal(err)
		}
		assertInts(t, "student", intIDs(got, func(s *student.Student) int64 { return s.ID() }),
			[]int64{101, 102})
	})

	t.Run("teacher", func(t *testing.T) {
		got, err := data.PageTeachers(ctx, 1, 10)
		if err != nil {
			t.Fatal(err)
		}
		assertInts(t, "teacher", intIDs(got, func(x *teacher.Teacher) int64 { return x.ID() }),
			[]int64{1, 2, 3})
	})

	t.Run("course_slot", func(t *testing.T) {
		got, err := data.PageCourseSlots(ctx, 1, 10)
		if err != nil {
			t.Fatal(err)
		}
		assertStrings(t, "course_slot", strIDs(got, func(s *courseSlot.CourseSlot) string { return s.ID() }),
			[]string{
				"550e8400-e29b-41d4-a716-446655440001",
				"550e8400-e29b-41d4-a716-446655440002",
				"550e8400-e29b-41d4-a716-446655440003",
				"550e8400-e29b-41d4-a716-446655440004",
				"550e8400-e29b-41d4-a716-446655440005",
				"550e8400-e29b-41d4-a716-446655440006",
			})
	})

	t.Run("course_slot_change", func(t *testing.T) {
		got, err := data.PageCourseSlotChanges(ctx, 1, 10)
		if err != nil {
			t.Fatal(err)
		}
		assertInts(t, "course_slot_change",
			intIDs(got, func(c *courseSlotChange.CourseSlotChange) int64 { return c.ID() }),
			[]int64{1, 2})
	})

	t.Run("course_enrollment", func(t *testing.T) {
		got, err := data.PageEnrollments(ctx, 1, 2)
		if err != nil {
			t.Fatal(err)
		}
		assertInts(t, "course_enrollment",
			intIDs(got, func(e *enrollment.CourseEnrollment) int64 { return e.ID() }),
			[]int64{1, 2})
	})

	t.Run("absence_record", func(t *testing.T) {
		got, err := data.PageAbsences(ctx, 1, 10)
		if err != nil {
			t.Fatal(err)
		}
		assertInts(t, "absence_record",
			intIDs(got, func(a *absence.AbsenceRecord) int64 { return a.ID() }), []int64{1, 2})
	})

	t.Run("student_makeup", func(t *testing.T) {
		got, err := data.PageMakeups(ctx, 1, 10)
		if err != nil {
			t.Fatal(err)
		}
		assertInts(t, "student_makeup",
			intIDs(got, func(m *makeup.StudentMakeup) int64 { return m.ID() }), []int64{1, 2})
	})

	t.Run("qualification", func(t *testing.T) {
		got, err := data.PageQualifications(ctx, 1, 2)
		if err != nil {
			t.Fatal(err)
		}
		assertInts(t, "qualification",
			intIDs(got, func(q *qualification.Qualification) int64 { return q.ID() }), []int64{1, 2})
	})
}

// --- 按外键列取列表 ---

func TestIntegrationListByID(t *testing.T) {
	data := newIntegrationData(t)
	ctx := context.Background()

	t.Run("course_slot 按课程", func(t *testing.T) {
		got, err := data.ListCourseSlotsByCourseID(ctx, "C001")
		if err != nil {
			t.Fatal(err)
		}
		assertStrings(t, "C001", strIDs(got, func(s *courseSlot.CourseSlot) string { return s.ID() }),
			[]string{
				"550e8400-e29b-41d4-a716-446655440001",
				"550e8400-e29b-41d4-a716-446655440002",
			})

		empty, err := data.ListCourseSlotsByCourseID(ctx, "C999")
		if err != nil {
			t.Fatal(err)
		}
		if empty == nil || len(empty) != 0 {
			t.Errorf("不存在的课程应返回非 nil 空切片，得到 %v", empty)
		}
	})

	t.Run("course_slot_change 按课程", func(t *testing.T) {
		got, err := data.ListCourseSlotChangesByCourseID(ctx, "C001")
		if err != nil {
			t.Fatal(err)
		}
		assertInts(t, "C001", intIDs(got, func(c *courseSlotChange.CourseSlotChange) int64 { return c.ID() }),
			[]int64{1, 2})

		got, err = data.ListCourseSlotChangesByCourseID(ctx, "C002")
		if err != nil {
			t.Fatal(err)
		}
		assertInts(t, "C002", intIDs(got, func(c *courseSlotChange.CourseSlotChange) int64 { return c.ID() }),
			[]int64{})
	})

	t.Run("absence 按学员 / 课程", func(t *testing.T) {
		byStudent, err := data.ListByStudentID(ctx, 101)
		if err != nil {
			t.Fatal(err)
		}
		assertInts(t, "student 101",
			intIDs(byStudent, func(a *absence.AbsenceRecord) int64 { return a.ID() }), []int64{1})

		byCourse, err := data.ListAbsencesByCourseID(ctx, "C002")
		if err != nil {
			t.Fatal(err)
		}
		assertInts(t, "course C002",
			intIDs(byCourse, func(a *absence.AbsenceRecord) int64 { return a.ID() }), []int64{2})
	})
}

// --- 关联查询：去重 + 按首次出现排序 ---

func TestIntegrationEnrollment(t *testing.T) {
	data := newIntegrationData(t)
	ctx := context.Background()
	q := data

	// 报名：101→C001、102→C002、103→C001(在读)、103→C002(已结业，不过滤)
	got, err := q.EnrolledCoursesByStudentID(ctx, 103)
	if err != nil {
		t.Fatal(err)
	}
	assertStrings(t, "103 上过的课", courseIDs(got), []string{"C001", "C002"})

	students, err := q.EnrolledStudentsByCourseID(ctx, "C001")
	if err != nil {
		t.Fatal(err)
	}
	assertInts(t, "C001 的学员", intIDs(students, func(s *student.Student) int64 { return s.ID() }),
		[]int64{101, 103})

	empty, err := q.EnrolledCoursesByStudentID(ctx, 999)
	if err != nil {
		t.Fatal(err)
	}
	assertStrings(t, "不存在的学员", courseIDs(empty), []string{})
}

func TestIntegrationMakeup(t *testing.T) {
	data := newIntegrationData(t)
	ctx := context.Background()
	q := data

	got, err := q.MakeupCoursesByStudentID(ctx, 101)
	if err != nil {
		t.Fatal(err)
	}
	assertStrings(t, "101 补过的课", courseIDs(got), []string{"C001"})

	students, err := q.MakeupStudentsByCourseID(ctx, "C001")
	if err != nil {
		t.Fatal(err)
	}
	assertInts(t, "C001 补课的学员", intIDs(students, func(s *student.Student) int64 { return s.ID() }),
		[]int64{101})

	// makeup 2 是「已预约未核销」，按内存版语义也要算「有补课记录」。
	got, err = q.MakeupCoursesByStudentID(ctx, 102)
	if err != nil {
		t.Fatal(err)
	}
	assertStrings(t, "102（已预约未补）", courseIDs(got), []string{"C002"})
}

func TestIntegrationQualification(t *testing.T) {
	data := newIntegrationData(t)
	ctx := context.Background()
	q := data

	// 资质 1→ct-0001、4→ct-0002；3→ct-0002(3)、ct-0001(5)
	types, err := q.CourseTypesByTeacherID(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	assertStrings(t, "讲师 1 的类型", strIDs(types, func(ct *courseType.CourseType) string { return ct.ID() }),
		[]string{"ct-0001", "ct-0002"})

	types, err = q.CourseTypesByTeacherID(ctx, 3)
	if err != nil {
		t.Fatal(err)
	}
	assertStrings(t, "讲师 3 的类型（按首次认证）",
		strIDs(types, func(ct *courseType.CourseType) string { return ct.ID() }),
		[]string{"ct-0002", "ct-0001"})

	teachers, err := q.TeachersByCourseTypeID(ctx, "ct-0001")
	if err != nil {
		t.Fatal(err)
	}
	assertInts(t, "ct-0001 的讲师", intIDs(teachers, func(x *teacher.Teacher) int64 { return x.ID() }),
		[]int64{1, 2, 3})

	teachers, err = q.TeachersByCourseTypeID(ctx, "ct-0002")
	if err != nil {
		t.Fatal(err)
	}
	assertInts(t, "ct-0002 的讲师", intIDs(teachers, func(x *teacher.Teacher) int64 { return x.ID() }),
		[]int64{3, 1})
}

func TestIntegrationIsTeacherQualified(t *testing.T) {
	data := newIntegrationData(t)
	ctx := context.Background()
	q := data

	cases := []struct {
		teacherID int64
		courseID  string
		want      bool
	}{
		{1, "C001", true},  // 讲师 1 有 ct-0001
		{2, "C001", true},  // 讲师 2 有 ct-0001
		{3, "C001", true},  // 讲师 3 的 ct-0001 资质已被吊销，但不按状态过滤
		{1, "C003", true},  // C003 属 ct-0002，讲师 1 有
		{2, "C003", false}, // 讲师 2 只有 ct-0001
		{1, "C999", false}, // 课程不存在
	}
	for _, tc := range cases {
		got, err := q.IsTeacherQualifiedForCourse(ctx, tc.teacherID, tc.courseID)
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.want {
			t.Errorf("IsTeacherQualifiedForCourse(%d, %s) = %v, want %v",
				tc.teacherID, tc.courseID, got, tc.want)
		}
	}
}
