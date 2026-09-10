package query

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

// classroomQuery 是 repoquery.ClassroomQuery 的内存实现。
type classroomQuery struct {
	data *memory.Data
}

// NewClassroomQuery 创建内存版教室查询。
func NewClassroomQuery(d *memory.Data) repoquery.ClassroomQuery {
	return &classroomQuery{data: d}
}

// Page 分页查询教室列表。
func (q *classroomQuery) Page(_ context.Context, page, pageSize int) ([]*classroom.Classroom, error) {
	return paginate(q.data.Classrooms(), page, pageSize), nil
}
