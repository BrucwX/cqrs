package imp

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/query/mysql"
	"cqrs/internal/core/course_scheduling/adapters/query/mysql/implement/help"
	"cqrs/internal/core/course_scheduling/adapters/query/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

// courseSlotQuery 是 repoquery.CourseSlotQuery 的 MySQL 实现。
type courseSlotQuery struct {
	data *mysql.Data
}

// 编译期断言：实现必须满足接口。
var _ repoquery.CourseSlotQuery = (*courseSlotQuery)(nil)

// NewCourseSlotQuery 创建 MySQL 版课表槽位查询。
func NewCourseSlotQuery(d *mysql.Data) repoquery.CourseSlotQuery {
	return &courseSlotQuery{data: d}
}

// Page 分页查询课表槽位列表，按 ID 升序。
func (q *courseSlotQuery) Page(ctx context.Context, page, pageSize int) ([]*courseSlot.CourseSlot, error) {
	limit, offset := help.LimitOffset(page, pageSize)
	return help.QueryAll(ctx, q.data.Conn(ctx), `
SELECT `+help.CourseSlotColumns+`
  FROM course_slot
 ORDER BY id
 LIMIT ? OFFSET ?`, []any{limit, offset}, help.ScanCourseSlot, model.CourseSlotToDO)
}

// ListByCourseID 根据课程 ID 获取课表槽位列表。
func (q *courseSlotQuery) ListByCourseID(ctx context.Context, courseID string) ([]*courseSlot.CourseSlot, error) {
	return help.QueryAll(ctx, q.data.Conn(ctx), `
SELECT `+help.CourseSlotColumns+`
  FROM course_slot
 WHERE course_id = ?
 ORDER BY id`, []any{courseID}, help.ScanCourseSlot, model.CourseSlotToDO)
}
