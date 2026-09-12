package repoImp

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/query/mysql"
	"cqrs/internal/core/course_scheduling/adapters/query/mysql/implement/help"
	"cqrs/internal/core/course_scheduling/adapters/query/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

// courseSlotChangeQuery 是 repoquery.CourseSlotChangeQuery 的 MySQL 实现。
type courseSlotChangeQuery struct {
	data *mysql.Data
}

// 编译期断言：实现必须满足接口。
var _ repoquery.CourseSlotChangeQuery = (*courseSlotChangeQuery)(nil)

// NewCourseSlotChangeQuery 创建 MySQL 版课表变更查询。
func NewCourseSlotChangeQuery(d *mysql.Data) repoquery.CourseSlotChangeQuery {
	return &courseSlotChangeQuery{data: d}
}

// Page 分页查询课表变更列表，按 ID 升序。
func (q *courseSlotChangeQuery) Page(ctx context.Context, page, pageSize int) ([]*courseSlotChange.CourseSlotChange, error) {
	limit, offset := help.LimitOffset(page, pageSize)
	return help.QueryAll(ctx, q.data.Conn(ctx), `
SELECT `+help.CourseSlotChangeColumns+`
  FROM course_slot_change
 ORDER BY id
 LIMIT ? OFFSET ?`, []any{limit, offset}, help.ScanCourseSlotChange, model.CourseSlotChangeToDO)
}

// ListByCourseID 根据课程 ID 获取课表变更列表。
//
// 接口没带 ctx，只能兜一个 Background，因此这个查询无法被取消
// （读侧本来也不参与写侧事务，影响有限）。
func (q *courseSlotChangeQuery) ListByCourseID(courseID string) ([]*courseSlotChange.CourseSlotChange, error) {
	ctx := context.Background()
	return help.QueryAll(ctx, q.data.Conn(ctx), `
SELECT `+help.CourseSlotChangeColumns+`
  FROM course_slot_change
 WHERE course_id = ?
 ORDER BY id`, []any{courseID}, help.ScanCourseSlotChange, model.CourseSlotChangeToDO)
}
