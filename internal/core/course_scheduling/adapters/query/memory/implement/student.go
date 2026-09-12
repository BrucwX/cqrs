package implement

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/query/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

// studentQuery 是 repoquery.StudentQuery 的内存实现。
type studentQuery struct {
	data *memory.Data
}

// NewStudentQuery 创建内存版学员查询。
func NewStudentQuery(d *memory.Data) repoquery.StudentQuery {
	return &studentQuery{data: d}
}

// Page 分页查询学员列表。
func (q *studentQuery) Page(_ context.Context, page, pageSize int) ([]*student.Student, error) {
	return paginate(q.data.Students(), page, pageSize), nil
}
