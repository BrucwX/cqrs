package implement

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/test_memory/query"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

// teacherQuery 是 repoquery.TeacherQuery 的内存实现。
type teacherQuery struct {
	data *query.Data
}

// NewTeacherQuery 创建内存版讲师查询。
func NewTeacherQuery(d *query.Data) repoquery.TeacherQuery {
	return &teacherQuery{data: d}
}

// Page 分页查询讲师列表。
func (q *teacherQuery) Page(_ context.Context, page, pageSize int) ([]*teacher.Teacher, error) {
	return paginate(q.data.Teachers(), page, pageSize), nil
}
