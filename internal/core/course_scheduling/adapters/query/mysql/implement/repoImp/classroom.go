package repoImp

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/query/mysql"
	"cqrs/internal/core/course_scheduling/adapters/query/mysql/implement/help"
	"cqrs/internal/core/course_scheduling/adapters/query/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

// classroomQuery 是 repoquery.ClassroomQuery 的 MySQL 实现。
type classroomQuery struct {
	data *mysql.Data
}

// 编译期断言：实现必须满足接口。
var _ repoquery.ClassroomQuery = (*classroomQuery)(nil)

// NewClassroomQuery 创建 MySQL 版教室查询。
func NewClassroomQuery(d *mysql.Data) repoquery.ClassroomQuery {
	return &classroomQuery{data: d}
}

// Page 分页查询教室列表，按 ID 升序。
func (q *classroomQuery) Page(ctx context.Context, page, pageSize int) ([]*classroom.Classroom, error) {
	limit, offset := help.LimitOffset(page, pageSize)
	return help.QueryAll(ctx, q.data.Conn(ctx), `
SELECT `+help.ClassroomColumns+`
  FROM classroom
 ORDER BY id
 LIMIT ? OFFSET ?`, []any{limit, offset}, help.ScanClassroom, model.ClassroomToDO)
}
