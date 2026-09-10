package command

import "cqrs/internal/core/course_scheduling/domain/aggregate/makeup"

// StudentMakeupCommand 补课申请命令接口
type StudentMakeupCommand interface {
	// Save 保存补课申请（新增或更新）
	Save(m *makeup.StudentMakeup) error
	// Delete 删除补课申请
	Delete(id int64) error
}
