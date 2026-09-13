package implement

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/memImp4test/query"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

// courseSlotQuery 是 repoquery.CourseSlotQuery 的内存实现。
type courseSlotQuery struct {
	data *query.Data
}

// NewCourseSlotQuery 创建内存版课表槽位查询。
func NewCourseSlotQuery(d *query.Data) repoquery.CourseSlotQuery {
	return &courseSlotQuery{data: d}
}

// PageCourseSlots 分页查询课表槽位列表。
func (q *courseSlotQuery) PageCourseSlots(_ context.Context, page, pageSize int) ([]*courseSlot.CourseSlot, error) {
	return paginate(q.data.CourseSlots(), page, pageSize), nil
}

// ListCourseSlotsByCourseID 根据课程 ID 获取课表槽位列表。
func (q *courseSlotQuery) ListCourseSlotsByCourseID(_ context.Context, courseID string) ([]*courseSlot.CourseSlot, error) {
	out := make([]*courseSlot.CourseSlot, 0)
	for _, item := range q.data.CourseSlots() {
		if item.CourseID() == courseID {
			out = append(out, item)
		}
	}
	return out, nil
}
