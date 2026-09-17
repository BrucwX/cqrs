package mysql

import (
	"context"

	"courseScheduling/internal/adapters/query/data/mysql/help"
	"courseScheduling/internal/adapters/query/data/mysql/model"
	"courseScheduling/internal/domain/aggregate/teacher"
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
