package query

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
)

// CourseSlotChangeQuery 课表变更查询接口
type CourseSlotChangeQuery interface {
	// PageCourseSlotChanges 分页查询课表变更列表
	PageCourseSlotChanges(ctx context.Context, page int, pageSize int) ([]*courseSlotChange.CourseSlotChange, error)
	// ListCourseSlotChangesByCourseID 根据课程 ID 获取课表变更列表
	ListCourseSlotChangesByCourseID(ctx context.Context, courseID string) ([]*courseSlotChange.CourseSlotChange, error)
}
