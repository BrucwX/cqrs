package query

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
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

// CourseTypesByTeacherID 根据讲师 ID 获取其有资质的课程类型列表。
//
// 不按资质状态过滤（已吊销/已过期也会返回）；结果去重。
// 若只要「当前有效」的资质，改用 item.IsEligible(time.Now()) == nil 判断。
func (q *qualificationQuery) CourseTypesByTeacherID(_ context.Context, teacherID int64) ([]*courseType.CourseType, error) {
	byID := indexBy(q.data.CourseTypes(), func(ct *courseType.CourseType) string { return ct.ID() })

	out := make([]*courseType.CourseType, 0)
	seen := make(map[string]struct{})
	for _, item := range q.data.Qualifications() {
		if item.TeacherID() != teacherID {
			continue
		}
		if _, ok := seen[item.CourseTypeID()]; ok {
			continue
		}
		seen[item.CourseTypeID()] = struct{}{}
		if ct, ok := byID[item.CourseTypeID()]; ok {
			out = append(out, ct)
		}
	}
	return out, nil
}

// TeachersByCourseTypeID 根据课程类型 ID 获取有授课资质的讲师列表。
//
// 同样不按资质状态过滤；结果去重。
func (q *qualificationQuery) TeachersByCourseTypeID(_ context.Context, courseTypeID string) ([]*teacher.Teacher, error) {
	byID := indexBy(q.data.Teachers(), func(t *teacher.Teacher) int64 { return t.ID() })

	out := make([]*teacher.Teacher, 0)
	seen := make(map[int64]struct{})
	for _, item := range q.data.Qualifications() {
		if item.CourseTypeID() != courseTypeID {
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

// IsTeacherQualifiedForCourse 判断讲师是否有教指定课程的资质。
//
// 先把课程换回它所属的课程类型，再看该讲师有没有这个类型的资质。
// 与其它资质查询一致：不按资质状态过滤（已吊销/已过期也算「有」）。
func (q *qualificationQuery) IsTeacherQualifiedForCourse(ctx context.Context, teacherID int64, courseID string) (bool, error) {
	courses := indexBy(q.data.Courses(), func(c *course.Course) string { return c.ID() })
	item, ok := courses[courseID]
	if !ok {
		return false, nil // 课程不存在，谈不上有资质
	}

	types, err := q.CourseTypesByTeacherID(ctx, teacherID)
	if err != nil {
		return false, err
	}
	for _, ct := range types {
		if ct.ID() == item.CourseTypeID() {
			return true, nil
		}
	}
	return false, nil
}
