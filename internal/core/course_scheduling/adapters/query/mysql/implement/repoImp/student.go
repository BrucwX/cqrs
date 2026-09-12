package repoImp

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/query/mysql"
	"cqrs/internal/core/course_scheduling/adapters/query/mysql/implement/help"
	"cqrs/internal/core/course_scheduling/adapters/query/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

// studentQuery 是 repoquery.StudentQuery 的 MySQL 实现。
type studentQuery struct {
	data *mysql.Data
}

// 编译期断言：实现必须满足接口。
var _ repoquery.StudentQuery = (*studentQuery)(nil)

// NewStudentQuery 创建 MySQL 版学员查询。
func NewStudentQuery(d *mysql.Data) repoquery.StudentQuery {
	return &studentQuery{data: d}
}

// Page 分页查询学员列表，按 ID 升序。
func (q *studentQuery) Page(ctx context.Context, page, pageSize int) ([]*student.Student, error) {
	limit, offset := help.LimitOffset(page, pageSize)
	return help.QueryAll(ctx, q.data.Conn(ctx), `
SELECT `+help.StudentColumns+`
  FROM student
 ORDER BY id
 LIMIT ? OFFSET ?`, []any{limit, offset}, help.ScanStudent, model.StudentToDO)
}
