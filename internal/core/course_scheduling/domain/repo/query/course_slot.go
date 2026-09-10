package query

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
)

// CourseSlotQuery 课表槽位查询接口
type CourseSlotQuery interface {
	// Page 分页查询课表槽位列表
	Page(ctx context.Context, page int, pageSize int) ([]*courseSlot.CourseSlot, error)
	// ListByCourseID 根据课程 ID 获取课表槽位列表
	ListByCourseID(courseID string) ([]*courseSlot.CourseSlot, error)
}
