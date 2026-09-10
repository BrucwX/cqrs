package query

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

// courseSlotQuery 是 repoquery.CourseSlotQuery 的内存实现。
type courseSlotQuery struct {
	data *memory.Data
}

// NewCourseSlotQuery 创建内存版课表槽位查询。
func NewCourseSlotQuery(d *memory.Data) repoquery.CourseSlotQuery {
	return &courseSlotQuery{data: d}
}

// Page 分页查询课表槽位列表。
func (q *courseSlotQuery) Page(_ context.Context, page, pageSize int) ([]*courseSlot.CourseSlot, error) {
	return paginate(q.data.CourseSlots(), page, pageSize), nil
}

// ListByCourseID 根据课程 ID 获取课表槽位列表。
func (q *courseSlotQuery) ListByCourseID(courseID string) ([]*courseSlot.CourseSlot, error) {
	out := make([]*courseSlot.CourseSlot, 0)
	for _, item := range q.data.CourseSlots() {
		if item.CourseID() == courseID {
			out = append(out, item)
		}
	}
	return out, nil
}
