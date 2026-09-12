package imp

import (
	"context"
	"testing"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
)

// 种子排期（星期几 / 时段 / 讲师 / 教室）：
//
//	C001 slot1 周一 09:00-11:00 讲师1 R101
//	C001 slot2 周三 09:00-11:00 讲师1 R101
//	C002 slot3 周一 14:00-16:00 讲师2 R102
//	C003 slot4 周三 09:00-11:00 讲师3 R102
//	C003 slot5 周五 09:00-11:00 讲师3 R102
//	C004 slot6 周二 09:00-11:00 讲师3 R102
//
// 报名：101→C001(在读)、102→C002(在读)、103→C001(在读)、103→C002(已结业)
//
// 冲突判定 = 同一星期几 且 日内时间区间重叠，与内存版一致。
func TestIntegrationAvailable(t *testing.T) {
	data := newIntegrationData(t)
	ctx := context.Background()
	q := NewCourseQuery(data)

	t.Run("学员", func(t *testing.T) {
		// 101 在学 C001（周一、周三 09-11）→ 撞上 C001 自己，以及同样占周三 09-11 的 C003
		got, err := q.AvailableForStudent(ctx, 101)
		if err != nil {
			t.Fatal(err)
		}
		assertStrings(t, "101", courseIDs(got), []string{"C002", "C004"})

		// 102 在学 C002（周一 14-16）→ 只撞 C002 自己
		got, err = q.AvailableForStudent(ctx, 102)
		if err != nil {
			t.Fatal(err)
		}
		assertStrings(t, "102", courseIDs(got), []string{"C001", "C003", "C004"})

		// 103 只有 C001 是在读（C002 已结业，不算占用）→ 同 101
		got, err = q.AvailableForStudent(ctx, 103)
		if err != nil {
			t.Fatal(err)
		}
		assertStrings(t, "103", courseIDs(got), []string{"C002", "C004"})

		got, err = q.AvailableForStudent(ctx, 999)
		if err != nil {
			t.Fatal(err)
		}
		assertStrings(t, "没有在学课程的学员", courseIDs(got), []string{"C001", "C002", "C003", "C004"})
	})

	t.Run("讲师", func(t *testing.T) {
		got, err := q.AvailableForTeacher(ctx, 1)
		if err != nil {
			t.Fatal(err)
		}
		assertStrings(t, "讲师 1", courseIDs(got), []string{"C002", "C004"})

		got, err = q.AvailableForTeacher(ctx, 2)
		if err != nil {
			t.Fatal(err)
		}
		assertStrings(t, "讲师 2", courseIDs(got), []string{"C001", "C003", "C004"})

		got, err = q.AvailableForTeacher(ctx, 3)
		if err != nil {
			t.Fatal(err)
		}
		// 讲师 3 占周三 / 周五 / 周二 09-11，C001 的周三 09-11 正好撞上
		assertStrings(t, "讲师 3", courseIDs(got), []string{"C002"})
	})

	t.Run("教室", func(t *testing.T) {
		got, err := q.AvailableForClassroom(ctx, "R101")
		if err != nil {
			t.Fatal(err)
		}
		assertStrings(t, "R101", courseIDs(got), []string{"C002", "C004"})

		// R102 被 C002/C003/C004 占满，C001 又和 C003 撞周三 09-11
		got, err = q.AvailableForClassroom(ctx, "R102")
		if err != nil {
			t.Fatal(err)
		}
		assertStrings(t, "R102", courseIDs(got), []string{})

		got, err = q.AvailableForClassroom(ctx, "R999")
		if err != nil {
			t.Fatal(err)
		}
		assertStrings(t, "没人用过的教室", courseIDs(got), []string{"C001", "C002", "C003", "C004"})
	})
}

// TestIntegrationFieldFidelity 盯一眼「列 -> PO -> DO」这条链路有没有走样：
// 值对象打平、TIME 列、NULL 列、枚举都从这里过。
func TestIntegrationFieldFidelity(t *testing.T) {
	data := newIntegrationData(t)
	ctx := context.Background()

	t.Run("classroom 的 Location 与容量", func(t *testing.T) {
		got, err := NewClassroomQuery(data).Page(ctx, 1, 1)
		if err != nil {
			t.Fatal(err)
		}
		if got[0].Location().FullName() != "A座-3F-301" {
			t.Errorf("Location = %q, want %q", got[0].Location().FullName(), "A座-3F-301")
		}
		if got[0].Capacity() != 30 {
			t.Errorf("Capacity = %d, want 30", got[0].Capacity())
		}
	})

	t.Run("course 的值对象", func(t *testing.T) {
		got, err := NewCourseQuery(data).Page(ctx, 1, 1)
		if err != nil {
			t.Fatal(err)
		}
		c := got[0]
		if c.Capacity().Max() != 30 || c.Capacity().Enrolled() != 2 {
			t.Errorf("Capacity = %d/%d, want 2/30", c.Capacity().Enrolled(), c.Capacity().Max())
		}
		if c.CourseTypeID() != "ct-0001" {
			t.Errorf("CourseTypeID = %q, want ct-0001", c.CourseTypeID())
		}
		if c.Period().TotalHours() == 0 {
			t.Error("Period.TotalHours 不应为 0")
		}
		if c.Enrollment().DropDeadline().IsZero() {
			t.Error("Enrollment.DropDeadline 不应为零值")
		}
	})

	t.Run("course_slot 的 TIME 列与星期", func(t *testing.T) {
		got, err := NewCourseSlotQuery(data).ListByCourseID(ctx, "C001")
		if err != nil {
			t.Fatal(err)
		}
		if got[0].TimeRange().Start().String() != "09:00" {
			t.Errorf("Start = %q, want %q", got[0].TimeRange().Start().String(), "09:00")
		}
		if got[0].TimeRange().End().String() != "11:00" {
			t.Errorf("End = %q, want %q", got[0].TimeRange().End().String(), "11:00")
		}
		if got[0].Weekday().String() != "Monday" {
			t.Errorf("Weekday = %v, want Monday", got[0].Weekday())
		}
		if got[0].ClassroomID() != "R101" {
			t.Errorf("ClassroomID = %q, want R101", got[0].ClassroomID())
		}
	})

	t.Run("enrollment 的可空列与枚举", func(t *testing.T) {
		got, err := NewCourseEnrollmentQuery(data).Page(ctx, 1, 4)
		if err != nil {
			t.Fatal(err)
		}
		if !got[0].CompletedAt().IsZero() {
			t.Errorf("completed_at 是 NULL，DO 里应是零值，得到 %v", got[0].CompletedAt())
		}
		if got[0].Status() != enrollment.StatusEnrolled {
			t.Errorf("Status = %v, want %v", got[0].Status(), enrollment.StatusEnrolled)
		}
		if got[3].CompletedAt().IsZero() {
			t.Error("报名 4 已结业，CompletedAt 不应为零值")
		}
	})

	t.Run("course_slot_change 的字符串快照与枚举", func(t *testing.T) {
		got, err := NewCourseSlotChangeQuery(data).ListByCourseID(ctx, "C001")
		if err != nil {
			t.Fatal(err)
		}
		if got[0].OriginalPlan().StartTimeStr() == "" {
			t.Error("OriginalPlan.StartTimeStr 不应为空")
		}
		if got[0].ChangeType() != courseSlotChange.TypeReschedule {
			t.Errorf("ChangeType = %v, want %v", got[0].ChangeType(), courseSlotChange.TypeReschedule)
		}
		if got[1].ChangeType() != courseSlotChange.TypeSubstitute {
			t.Errorf("ChangeType = %v, want %v", got[1].ChangeType(), courseSlotChange.TypeSubstitute)
		}
	})
}
