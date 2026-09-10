package courseSlot

import "cqrs/internal/core/course_scheduling/domain/repo/query"

// Handler 课表槽位查询处理器
type Handler struct {
	Query query.CourseSlotQuery
}
