package implement

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/query/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

// courseSlotChangeQuery 是 repoquery.CourseSlotChangeQuery 的内存实现。
type courseSlotChangeQuery struct {
	data *memory.Data
}

// NewCourseSlotChangeQuery 创建内存版课表变更查询。
func NewCourseSlotChangeQuery(d *memory.Data) repoquery.CourseSlotChangeQuery {
	return &courseSlotChangeQuery{data: d}
}

// Page 分页查询课表变更列表。
func (q *courseSlotChangeQuery) Page(_ context.Context, page, pageSize int) ([]*courseSlotChange.CourseSlotChange, error) {
	return paginate(q.data.CourseSlotChanges(), page, pageSize), nil
}

// ListByCourseID 根据课程 ID 获取课表变更列表。
func (q *courseSlotChangeQuery) ListByCourseID(courseID string) ([]*courseSlotChange.CourseSlotChange, error) {
	out := make([]*courseSlotChange.CourseSlotChange, 0)
	for _, item := range q.data.CourseSlotChanges() {
		if item.CourseID() == courseID {
			out = append(out, item)
		}
	}
	return out, nil
}
