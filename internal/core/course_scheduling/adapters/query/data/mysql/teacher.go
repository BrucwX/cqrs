package mysql

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/query/data/mysql/help"
	"cqrs/internal/core/course_scheduling/adapters/query/data/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
)

// Page 分页查询讲师列表，按 ID 升序。
func (q *Data) PageTeachers(ctx context.Context, page, pageSize int) ([]*teacher.Teacher, error) {
	limit, offset := help.LimitOffset(page, pageSize)
	return help.QueryAll(ctx, q.Conn(ctx), `
SELECT `+help.TeacherColumns+`
  FROM teacher
 ORDER BY id
 LIMIT ? OFFSET ?`, []any{limit, offset}, help.ScanTeacher, model.TeacherToDO)
}
