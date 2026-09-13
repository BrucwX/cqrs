package mysql

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/query/data/mysql/help"
	"cqrs/internal/core/course_scheduling/adapters/query/data/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
)

// Page 分页查询教室列表，按 ID 升序。
func (q *Data) PageClassrooms(ctx context.Context, page, pageSize int) ([]*classroom.Classroom, error) {
	limit, offset := help.LimitOffset(page, pageSize)
	return help.QueryAll(ctx, q.Conn(ctx), `
SELECT `+help.ClassroomColumns+`
  FROM classroom
 ORDER BY id
 LIMIT ? OFFSET ?`, []any{limit, offset}, help.ScanClassroom, model.ClassroomToDO)
}
