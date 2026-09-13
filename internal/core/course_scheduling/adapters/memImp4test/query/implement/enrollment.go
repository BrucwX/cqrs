package implement

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/memImp4test/query"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

// courseEnrollmentQuery 是 repoquery.CourseEnrollmentQuery 的内存实现。
type courseEnrollmentQuery struct {
	data *query.Data
}

// NewCourseEnrollmentQuery 创建内存版课程注册查询。
func NewCourseEnrollmentQuery(d *query.Data) repoquery.CourseEnrollmentQuery {
	return &courseEnrollmentQuery{data: d}
}

// PageEnrollments 分页查询课程注册列表。
func (q *courseEnrollmentQuery) PageEnrollments(_ context.Context, page, pageSize int) ([]*enrollment.CourseEnrollment, error) {
	return paginate(q.data.Enrollments(), page, pageSize), nil
}

// EnrolledCoursesByStudentID 根据学员 ID 获取其注册过的课程列表。
//
// 不按注册状态过滤（已完成/已退课的课程也会返回）；结果去重，顺序按注册时间先后。
// 若只想取在学课程，把 IsActive 判断加回来即可。
func (q *courseEnrollmentQuery) EnrolledCoursesByStudentID(_ context.Context, studentID int64) ([]*course.Course, error) {
	byID := indexBy(q.data.Courses(), func(c *course.Course) string { return c.ID() })

	out := make([]*course.Course, 0)
	seen := make(map[string]struct{})
	for _, e := range q.data.Enrollments() {
		if e.StudentID() != studentID {
			continue
		}
		if _, ok := seen[e.CourseID()]; ok {
			continue
		}
		seen[e.CourseID()] = struct{}{}
		if c, ok := byID[e.CourseID()]; ok {
			out = append(out, c)
		}
	}
	return out, nil
}

// EnrolledStudentsByCourseID 根据课程 ID 获取注册过的学员列表。
//
// 同样不按注册状态过滤；结果去重，顺序按注册时间先后。
func (q *courseEnrollmentQuery) EnrolledStudentsByCourseID(_ context.Context, courseID string) ([]*student.Student, error) {
	byID := indexBy(q.data.Students(), func(s *student.Student) int64 { return s.ID() })

	out := make([]*student.Student, 0)
	seen := make(map[int64]struct{})
	for _, e := range q.data.Enrollments() {
		if e.CourseID() != courseID {
			continue
		}
		if _, ok := seen[e.StudentID()]; ok {
			continue
		}
		seen[e.StudentID()] = struct{}{}
		if s, ok := byID[e.StudentID()]; ok {
			out = append(out, s)
		}
	}
	return out, nil
}
