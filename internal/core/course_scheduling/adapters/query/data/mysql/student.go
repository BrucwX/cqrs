package mysql

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/query/data/mysql/help"
	"cqrs/internal/core/course_scheduling/adapters/query/data/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
)

// Page 分页查询学员列表，按 ID 升序。
func (q *Data) PageStudents(ctx context.Context, page, pageSize int) ([]*student.Student, error) {
	limit, offset := help.LimitOffset(page, pageSize)
	return help.QueryAll(ctx, q.Conn(ctx), `
SELECT `+help.StudentColumns+`
  FROM student
 ORDER BY id
 LIMIT ? OFFSET ?`, []any{limit, offset}, help.ScanStudent, model.StudentToDO)
}
