package query

import (
	"context"
	"time"

	"cqrs/internal/core/course_scheduling/adapters/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

// courseQuery 是 repoquery.CourseQuery 的内存实现。
type courseQuery struct {
	data *memory.Data
}

// NewCourseQuery 创建内存版课程查询。
func NewCourseQuery(d *memory.Data) repoquery.CourseQuery {
	return &courseQuery{data: d}
}

// Page 分页查询课程列表。
func (q *courseQuery) Page(_ context.Context, page, pageSize int) ([]*course.Course, error) {
	return paginate(q.data.Courses(), page, pageSize), nil
}

// AvailableForStudent 返回与学员「当前在学课程」不冲突的课程。
//
// 占用集 = 该学员所有生效中（IsActive）注册记录对应课程的全部周排期槽位。
func (q *courseQuery) AvailableForStudent(_ context.Context, studentID int64) ([]*course.Course, error) {
	enrolled := make(map[string]struct{})
	for _, e := range q.data.Enrollments() {
		if e.StudentID() == studentID && e.IsActive() {
			enrolled[e.CourseID()] = struct{}{}
		}
	}

	return q.available(func(slot *courseSlot.CourseSlot) bool {
		_, ok := enrolled[slot.CourseID()]
		return ok
	}), nil
}

// AvailableForTeacher 返回与讲师「现有排课」不冲突的课程。
//
// 占用集 = 该讲师已被指派（TeacherID 命中）的全部周排期槽位。
func (q *courseQuery) AvailableForTeacher(_ context.Context, teacherID int64) ([]*course.Course, error) {
	return q.available(func(slot *courseSlot.CourseSlot) bool {
		return slot.TeacherID() == teacherID
	}), nil
}

// AvailableForClassroom 返回与教室「现有排课」不冲突的课程。
//
// 占用集 = 该教室已被占用（ClassroomID 命中）的全部周排期槽位。
func (q *courseQuery) AvailableForClassroom(_ context.Context, classroomID string) ([]*course.Course, error) {
	return q.available(func(slot *courseSlot.CourseSlot) bool {
		return slot.ClassroomID() == classroomID
	}), nil
}

// --- 内部实现 ---

// busySlot 是一段已被占用的周内固定时间（星期几 + 日内时间区间）。
type busySlot struct {
	weekday time.Weekday
	span    courseSlot.DayTimeRange
}

// slotPredicate 判定某个槽位是否属于「已占用」集合。
type slotPredicate func(*courseSlot.CourseSlot) bool

// available 返回所有槽位都不与「已占用」时间冲突的课程。
//
// 没有排期槽位的课程不会被判定为冲突，因此会出现在结果里。
func (q *courseQuery) available(occupied slotPredicate) []*course.Course {
	busy := make([]busySlot, 0)
	for _, slot := range q.data.CourseSlots() {
		if occupied(slot) {
			busy = append(busy, busySlot{weekday: slot.Weekday(), span: slot.TimeRange()})
		}
	}

	slotsByCourse := make(map[string][]*courseSlot.CourseSlot)
	for _, slot := range q.data.CourseSlots() {
		slotsByCourse[slot.CourseID()] = append(slotsByCourse[slot.CourseID()], slot)
	}

	courses := q.data.Courses()
	out := make([]*course.Course, 0, len(courses))
	for _, c := range courses {
		if !conflictsWithAny(slotsByCourse[c.ID()], busy) {
			out = append(out, c)
		}
	}
	return out
}

// conflictsWithAny 判断候选课程是否与已占用时间冲突。
//
// 冲突判定 = 同一星期几 且 日内时间区间重叠。
// 这里刻意不看讲师/教室：对学员而言，两门课时间撞了就是上不了。
func conflictsWithAny(slots []*courseSlot.CourseSlot, busy []busySlot) bool {
	for _, slot := range slots {
		for _, b := range busy {
			if slot.Weekday() == b.weekday && slot.TimeRange().Overlaps(b.span) {
				return true
			}
		}
	}
	return false
}
