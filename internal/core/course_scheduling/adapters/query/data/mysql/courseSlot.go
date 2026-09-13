package mysql

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/query/data/mysql/help"
	"cqrs/internal/core/course_scheduling/adapters/query/data/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
)

// Page 分页查询课表槽位列表，按 ID 升序。
func (q *Data) PageCourseSlots(ctx context.Context, page, pageSize int) ([]*courseSlot.CourseSlot, error) {
	limit, offset := help.LimitOffset(page, pageSize)
	return help.QueryAll(ctx, q.Conn(ctx), `
SELECT `+help.CourseSlotColumns+`
  FROM course_slot
 ORDER BY id
 LIMIT ? OFFSET ?`, []any{limit, offset}, help.ScanCourseSlot, model.CourseSlotToDO)
}

// ListByCourseID 根据课程 ID 获取课表槽位列表。
func (q *Data) ListCourseSlotsByCourseID(ctx context.Context, courseID string) ([]*courseSlot.CourseSlot, error) {
	return help.QueryAll(ctx, q.Conn(ctx), `
SELECT `+help.CourseSlotColumns+`
  FROM course_slot
 WHERE course_id = ?
 ORDER BY id`, []any{courseID}, help.ScanCourseSlot, model.CourseSlotToDO)
}
