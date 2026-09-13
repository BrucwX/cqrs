package mysql

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/query/data/mysql/help"
	"cqrs/internal/core/course_scheduling/adapters/query/data/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
)

// Page 分页查询授课资质列表，按 ID 升序。
func (q *Data) PageQualifications(ctx context.Context, page, pageSize int) ([]*qualification.Qualification, error) {
	limit, offset := help.LimitOffset(page, pageSize)
	return help.QueryAll(ctx, q.Conn(ctx), `
SELECT `+help.QualificationColumns+`
  FROM qualification
 ORDER BY id
 LIMIT ? OFFSET ?`, []any{limit, offset}, help.ScanQualification, model.QualificationToDO)
}

// CourseTypesByTeacherID 根据讲师 ID 获取其有资质的课程类型列表。
//
// 不按资质状态过滤（已吊销、已过期的也会返回）；去重，
// 顺序按第一条资质记录。
func (q *Data) CourseTypesByTeacherID(ctx context.Context, teacherID int64) ([]*courseType.CourseType, error) {
	query := `
SELECT ` + help.Qualify(help.CourseTypeColumns, "ct") + `
  FROM course_type ct
  JOIN (SELECT q.course_type_id, MIN(q.id) AS first_id
          FROM qualification q
         WHERE q.teacher_id = ?
         GROUP BY q.course_type_id) f ON f.course_type_id = ct.id
 ORDER BY f.first_id`

	return help.QueryAll(ctx, q.Conn(ctx), query, []any{teacherID}, help.ScanCourseType, model.CourseTypeToDO)
}

// TeachersByCourseTypeID 根据课程类型 ID 获取有授课资质的讲师列表。
//
// 同样不按资质状态过滤；去重，顺序按第一条资质记录。
func (q *Data) TeachersByCourseTypeID(ctx context.Context, courseTypeID string) ([]*teacher.Teacher, error) {
	query := `
SELECT ` + help.Qualify(help.TeacherColumns, "t") + `
  FROM teacher t
  JOIN (SELECT q.teacher_id, MIN(q.id) AS first_id
          FROM qualification q
         WHERE q.course_type_id = ?
         GROUP BY q.teacher_id) f ON f.teacher_id = t.id
 ORDER BY f.first_id`

	return help.QueryAll(ctx, q.Conn(ctx), query, []any{courseTypeID}, help.ScanTeacher, model.TeacherToDO)
}

// IsTeacherQualifiedForCourse 判断讲师是否有教指定课程的资质。
//
// 链路：课程 → 所属课程类型 → 该讲师是否有该类型的资质。
// 课程不存在时 EXISTS 自然为 false，与内存版一致。
// 与其它资质查询一致：不按资质状态过滤（已吊销、已过期的也算「有」）。
func (q *Data) IsTeacherQualifiedForCourse(ctx context.Context, teacherID int64, courseID string) (bool, error) {
	const query = `
SELECT EXISTS (
       SELECT 1
         FROM course c
         JOIN qualification q ON q.course_type_id = c.course_type_id
        WHERE c.id = ? AND q.teacher_id = ?)`

	return help.Exists(ctx, q.Conn(ctx), query, courseID, teacherID)
}
