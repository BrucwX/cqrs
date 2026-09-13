package mysql

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/query/data/mysql/help"
	"cqrs/internal/core/course_scheduling/adapters/query/data/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
)

// Page 分页查询补课申请列表，按 ID 升序。
func (q *Data) PageMakeups(ctx context.Context, page, pageSize int) ([]*makeup.StudentMakeup, error) {
	limit, offset := help.LimitOffset(page, pageSize)
	return help.QueryAll(ctx, q.Conn(ctx), `
SELECT `+help.MakeupColumns+`
  FROM student_makeup
 ORDER BY id
 LIMIT ? OFFSET ?`, []any{limit, offset}, help.ScanMakeup, model.MakeupToDO)
}

// MakeupCoursesByStudentID 根据学员 ID 获取其有补课记录的课程列表。
//
// 不按补课状态过滤（已预约未补、已取消也会返回）；结果去重，
// 顺序按第一条补课记录。
func (q *Data) MakeupCoursesByStudentID(ctx context.Context, studentID int64) ([]*course.Course, error) {
	query := `
SELECT ` + help.Qualify(help.CourseColumns, "c") + `
  FROM course c
  JOIN (SELECT m.course_id, MIN(m.id) AS first_id
          FROM student_makeup m
         WHERE m.student_id = ?
         GROUP BY m.course_id) f ON f.course_id = c.id
 ORDER BY f.first_id`

	return help.QueryAll(ctx, q.Conn(ctx), query, []any{studentID}, help.ScanCourse, model.CourseToDO)
}

// MakeupStudentsByCourseID 根据课程 ID 获取有补课记录的学员列表。
//
// 同样不按补课状态过滤；去重，顺序按第一条补课记录。
func (q *Data) MakeupStudentsByCourseID(ctx context.Context, courseID string) ([]*student.Student, error) {
	query := `
SELECT ` + help.Qualify(help.StudentColumns, "s") + `
  FROM student s
  JOIN (SELECT m.student_id, MIN(m.id) AS first_id
          FROM student_makeup m
         WHERE m.course_id = ?
         GROUP BY m.student_id) f ON f.student_id = s.id
 ORDER BY f.first_id`

	return help.QueryAll(ctx, q.Conn(ctx), query, []any{courseID}, help.ScanStudent, model.StudentToDO)
}
