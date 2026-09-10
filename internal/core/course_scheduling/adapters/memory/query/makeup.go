package query

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

// studentMakeupQuery 是 repoquery.StudentMakeupQuery 的内存实现。
type studentMakeupQuery struct {
	data *memory.Data
}

// NewStudentMakeupQuery 创建内存版补课申请查询。
func NewStudentMakeupQuery(d *memory.Data) repoquery.StudentMakeupQuery {
	return &studentMakeupQuery{data: d}
}

// Page 分页查询补课申请列表。
func (q *studentMakeupQuery) Page(_ context.Context, page, pageSize int) ([]*makeup.StudentMakeup, error) {
	return paginate(q.data.Makeups(), page, pageSize), nil
}

// CoursesByStudentID 根据学员 ID 获取其有补课记录的课程列表。
//
// 不按申请状态过滤（待审批/已驳回也会返回）；结果去重。
// 若只想统计「确实补上了」的课程，把 status == StatusCompleted 加回来即可。
func (q *studentMakeupQuery) CoursesByStudentID(_ context.Context, studentID int64) ([]*course.Course, error) {
	byID := indexBy(q.data.Courses(), func(c *course.Course) string { return c.ID() })

	out := make([]*course.Course, 0)
	seen := make(map[string]struct{})
	for _, m := range q.data.Makeups() {
		if m.StudentID() != studentID {
			continue
		}
		if _, ok := seen[m.CourseID()]; ok {
			continue
		}
		seen[m.CourseID()] = struct{}{}
		if c, ok := byID[m.CourseID()]; ok {
			out = append(out, c)
		}
	}
	return out, nil
}

// StudentsByCourseID 根据课程 ID 获取有补课记录的学员列表。
//
// 同样不按申请状态过滤；结果去重。
func (q *studentMakeupQuery) StudentsByCourseID(_ context.Context, courseID string) ([]*student.Student, error) {
	byID := indexBy(q.data.Students(), func(s *student.Student) int64 { return s.ID() })

	out := make([]*student.Student, 0)
	seen := make(map[int64]struct{})
	for _, m := range q.data.Makeups() {
		if m.CourseID() != courseID {
			continue
		}
		if _, ok := seen[m.StudentID()]; ok {
			continue
		}
		seen[m.StudentID()] = struct{}{}
		if s, ok := byID[m.StudentID()]; ok {
			out = append(out, s)
		}
	}
	return out, nil
}
