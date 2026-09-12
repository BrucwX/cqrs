package repoImp

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/query/mysql"
	"cqrs/internal/core/course_scheduling/adapters/query/mysql/implement/help"
	"cqrs/internal/core/course_scheduling/adapters/query/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

// courseEnrollmentQuery 是 repoquery.CourseEnrollmentQuery 的 MySQL 实现。
type courseEnrollmentQuery struct {
	data *mysql.Data
}

// 编译期断言：实现必须满足接口。
var _ repoquery.CourseEnrollmentQuery = (*courseEnrollmentQuery)(nil)

// NewCourseEnrollmentQuery 创建 MySQL 版课程注册查询。
func NewCourseEnrollmentQuery(d *mysql.Data) repoquery.CourseEnrollmentQuery {
	return &courseEnrollmentQuery{data: d}
}

// Page 分页查询课程注册列表，按 ID 升序。
func (q *courseEnrollmentQuery) Page(ctx context.Context, page, pageSize int) ([]*enrollment.CourseEnrollment, error) {
	limit, offset := help.LimitOffset(page, pageSize)
	return help.QueryAll(ctx, q.data.Conn(ctx), `
SELECT `+help.EnrollmentColumns+`
  FROM course_enrollment
 ORDER BY id
 LIMIT ? OFFSET ?`, []any{limit, offset}, help.ScanEnrollment, model.EnrollmentToDO)
}

// CoursesByStudentID 根据学员 ID 获取其注册过的课程列表。
//
// 不按注册状态过滤（已完成/已退课的课程也会返回）；结果去重，
// 顺序按「第一次报名」的先后，与内存版按报名记录 ID 升序一致。
//
// JOIN 派生表里取 MIN(id) 就是为了拿到这个「首次报名」的顺序：
// 直接 JOIN 会一个课程出多行，GROUP BY 又没法可靠地用它排序。
func (q *courseEnrollmentQuery) CoursesByStudentID(ctx context.Context, studentID int64) ([]*course.Course, error) {
	query := `
SELECT ` + help.Qualify(help.CourseColumns, "c") + `
  FROM course c
  JOIN (SELECT e.course_id, MIN(e.id) AS first_id
          FROM course_enrollment e
         WHERE e.student_id = ?
         GROUP BY e.course_id) f ON f.course_id = c.id
 ORDER BY f.first_id`

	return help.QueryAll(ctx, q.data.Conn(ctx), query, []any{studentID}, help.ScanCourse, model.CourseToDO)
}

// StudentsByCourseID 根据课程 ID 获取注册过的学员列表。
//
// 同样不按注册状态过滤；去重，顺序按「第一次报名」的先后。
func (q *courseEnrollmentQuery) StudentsByCourseID(ctx context.Context, courseID string) ([]*student.Student, error) {
	query := `
SELECT ` + help.Qualify(help.StudentColumns, "s") + `
  FROM student s
  JOIN (SELECT e.student_id, MIN(e.id) AS first_id
          FROM course_enrollment e
         WHERE e.course_id = ?
         GROUP BY e.student_id) f ON f.student_id = s.id
 ORDER BY f.first_id`

	return help.QueryAll(ctx, q.data.Conn(ctx), query, []any{courseID}, help.ScanStudent, model.StudentToDO)
}
