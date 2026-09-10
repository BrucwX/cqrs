package query

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

// qualificationQuery 是 repoquery.QualificationQuery 的内存实现。
type qualificationQuery struct {
	data *memory.Data
}

// NewQualificationQuery 创建内存版授课资质查询。
func NewQualificationQuery(d *memory.Data) repoquery.QualificationQuery {
	return &qualificationQuery{data: d}
}

// Page 分页查询授课资质列表。
func (q *qualificationQuery) Page(_ context.Context, page, pageSize int) ([]*qualification.Qualification, error) {
	return paginate(q.data.Qualifications(), page, pageSize), nil
}

// CoursesByTeacherID 根据讲师 ID 获取其有资质的课程列表。
//
// 不按资质状态过滤（已吊销/已过期也会返回）；结果去重。
// 若只要「当前有效」的资质，改用 q.IsEligible(time.Now()) == nil 判断。
func (q *qualificationQuery) CoursesByTeacherID(_ context.Context, teacherID int64) ([]*course.Course, error) {
	byID := indexBy(q.data.Courses(), func(c *course.Course) string { return c.ID() })

	out := make([]*course.Course, 0)
	seen := make(map[string]struct{})
	for _, item := range q.data.Qualifications() {
		if item.TeacherID() != teacherID {
			continue
		}
		if _, ok := seen[item.CourseID()]; ok {
			continue
		}
		seen[item.CourseID()] = struct{}{}
		if c, ok := byID[item.CourseID()]; ok {
			out = append(out, c)
		}
	}
	return out, nil
}

// TeachersByCourseID 根据课程 ID 获取有授课资质的讲师列表。
//
// 同样不按资质状态过滤；结果去重。
func (q *qualificationQuery) TeachersByCourseID(_ context.Context, courseID string) ([]*teacher.Teacher, error) {
	byID := indexBy(q.data.Teachers(), func(t *teacher.Teacher) int64 { return t.ID() })

	out := make([]*teacher.Teacher, 0)
	seen := make(map[int64]struct{})
	for _, item := range q.data.Qualifications() {
		if item.CourseID() != courseID {
			continue
		}
		if _, ok := seen[item.TeacherID()]; ok {
			continue
		}
		seen[item.TeacherID()] = struct{}{}
		if t, ok := byID[item.TeacherID()]; ok {
			out = append(out, t)
		}
	}
	return out, nil
}
