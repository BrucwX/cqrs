package mysql

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/query/data/mysql/help"
	"cqrs/internal/core/course_scheduling/adapters/query/data/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
)

// Page 分页查询课表变更列表，按 ID 升序。
func (q *Data) PageCourseSlotChanges(ctx context.Context, page, pageSize int) ([]*courseSlotChange.CourseSlotChange, error) {
	limit, offset := help.LimitOffset(page, pageSize)
	return help.QueryAll(ctx, q.Conn(ctx), `
SELECT `+help.CourseSlotChangeColumns+`
  FROM course_slot_change
 ORDER BY id
 LIMIT ? OFFSET ?`, []any{limit, offset}, help.ScanCourseSlotChange, model.CourseSlotChangeToDO)
}

// ListByCourseID 根据课程 ID 获取课表变更列表。
func (q *Data) ListCourseSlotChangesByCourseID(ctx context.Context, courseID string) ([]*courseSlotChange.CourseSlotChange, error) {
	return help.QueryAll(ctx, q.Conn(ctx), `
SELECT `+help.CourseSlotChangeColumns+`
  FROM course_slot_change
 WHERE course_id = ?
 ORDER BY id`, []any{courseID}, help.ScanCourseSlotChange, model.CourseSlotChangeToDO)
}
