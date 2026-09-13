package implement

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/memImp4test/query"
	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

// studentQuery 是 repoquery.StudentQuery 的内存实现。
type studentQuery struct {
	data *query.Data
}

// NewStudentQuery 创建内存版学员查询。
func NewStudentQuery(d *query.Data) repoquery.StudentQuery {
	return &studentQuery{data: d}
}

// PageStudents 分页查询学员列表。
func (q *studentQuery) PageStudents(_ context.Context, page, pageSize int) ([]*student.Student, error) {
	return paginate(q.data.Students(), page, pageSize), nil
}
