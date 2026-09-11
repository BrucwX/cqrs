package command

import (
	"errors"

	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
)

var (
	// ErrMakeupRequired 传入的补课预约为空。
	ErrMakeupRequired = errors.New("student makeup is required")
	// ErrMakeupNotFound 指定的补课预约不存在。
	ErrMakeupNotFound = errors.New("student makeup not found")
)

// StudentMakeupCommand 补课申请命令接口
type StudentMakeupCommand interface {
	// Save 保存补课申请（新增或更新）
	Save(m *makeup.StudentMakeup) error
	// Delete 删除补课申请
	Delete(id int64) error
}
