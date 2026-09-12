package repoImp

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/query/mysql"
	"cqrs/internal/core/course_scheduling/adapters/query/mysql/implement/help"
	"cqrs/internal/core/course_scheduling/adapters/query/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

// teacherQuery 是 repoquery.TeacherQuery 的 MySQL 实现。
type teacherQuery struct {
	data *mysql.Data
}

// 编译期断言：实现必须满足接口。
var _ repoquery.TeacherQuery = (*teacherQuery)(nil)

// NewTeacherQuery 创建 MySQL 版讲师查询。
func NewTeacherQuery(d *mysql.Data) repoquery.TeacherQuery {
	return &teacherQuery{data: d}
}

// Page 分页查询讲师列表，按 ID 升序。
func (q *teacherQuery) Page(ctx context.Context, page, pageSize int) ([]*teacher.Teacher, error) {
	limit, offset := help.LimitOffset(page, pageSize)
	return help.QueryAll(ctx, q.data.Conn(ctx), `
SELECT `+help.TeacherColumns+`
  FROM teacher
 ORDER BY id
 LIMIT ? OFFSET ?`, []any{limit, offset}, help.ScanTeacher, model.TeacherToDO)
}
