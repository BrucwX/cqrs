package implement

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/memImp4test/query"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

// courseSlotChangeQuery 是 repoquery.CourseSlotChangeQuery 的内存实现。
type courseSlotChangeQuery struct {
	data *query.Data
}

// NewCourseSlotChangeQuery 创建内存版课表变更查询。
func NewCourseSlotChangeQuery(d *query.Data) repoquery.CourseSlotChangeQuery {
	return &courseSlotChangeQuery{data: d}
}

// PageCourseSlotChanges 分页查询课表变更列表。
func (q *courseSlotChangeQuery) PageCourseSlotChanges(_ context.Context, page, pageSize int) ([]*courseSlotChange.CourseSlotChange, error) {
	return paginate(q.data.CourseSlotChanges(), page, pageSize), nil
}

// ListCourseSlotChangesByCourseID 根据课程 ID 获取课表变更列表。
func (q *courseSlotChangeQuery) ListCourseSlotChangesByCourseID(_ context.Context, courseID string) ([]*courseSlotChange.CourseSlotChange, error) {
	out := make([]*courseSlotChange.CourseSlotChange, 0)
	for _, item := range q.data.CourseSlotChanges() {
		if item.CourseID() == courseID {
			out = append(out, item)
		}
	}
	return out, nil
}
