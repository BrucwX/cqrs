package repoImp

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/query/mysql"
	"cqrs/internal/core/course_scheduling/adapters/query/mysql/implement/help"
	"cqrs/internal/core/course_scheduling/adapters/query/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

// studentMakeupQuery 是 repoquery.StudentMakeupQuery 的 MySQL 实现。
type studentMakeupQuery struct {
	data *mysql.Data
}

// 编译期断言：实现必须满足接口。
var _ repoquery.StudentMakeupQuery = (*studentMakeupQuery)(nil)

// NewStudentMakeupQuery 创建 MySQL 版补课申请查询。
func NewStudentMakeupQuery(d *mysql.Data) repoquery.StudentMakeupQuery {
	return &studentMakeupQuery{data: d}
}

// Page 分页查询补课申请列表，按 ID 升序。
func (q *studentMakeupQuery) Page(ctx context.Context, page, pageSize int) ([]*makeup.StudentMakeup, error) {
	limit, offset := help.LimitOffset(page, pageSize)
	return help.QueryAll(ctx, q.data.Conn(ctx), `
SELECT `+help.MakeupColumns+`
  FROM student_makeup
 ORDER BY id
 LIMIT ? OFFSET ?`, []any{limit, offset}, help.ScanMakeup, model.MakeupToDO)
}

// CoursesByStudentID 根据学员 ID 获取其有补课记录的课程列表。
//
// 不按补课状态过滤（已预约未补、已取消也会返回）；结果去重，
// 顺序按第一条补课记录。
func (q *studentMakeupQuery) CoursesByStudentID(ctx context.Context, studentID int64) ([]*course.Course, error) {
	query := `
SELECT ` + help.Qualify(help.CourseColumns, "c") + `
  FROM course c
  JOIN (SELECT m.course_id, MIN(m.id) AS first_id
          FROM student_makeup m
         WHERE m.student_id = ?
         GROUP BY m.course_id) f ON f.course_id = c.id
 ORDER BY f.first_id`

	return help.QueryAll(ctx, q.data.Conn(ctx), query, []any{studentID}, help.ScanCourse, model.CourseToDO)
}

// StudentsByCourseID 根据课程 ID 获取有补课记录的学员列表。
//
// 同样不按补课状态过滤；去重，顺序按第一条补课记录。
func (q *studentMakeupQuery) StudentsByCourseID(ctx context.Context, courseID string) ([]*student.Student, error) {
	query := `
SELECT ` + help.Qualify(help.StudentColumns, "s") + `
  FROM student s
  JOIN (SELECT m.student_id, MIN(m.id) AS first_id
          FROM student_makeup m
         WHERE m.course_id = ?
         GROUP BY m.student_id) f ON f.student_id = s.id
 ORDER BY f.first_id`

	return help.QueryAll(ctx, q.data.Conn(ctx), query, []any{courseID}, help.ScanStudent, model.StudentToDO)
}
