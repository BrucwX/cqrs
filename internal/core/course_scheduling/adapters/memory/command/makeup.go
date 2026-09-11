package command

import (
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// MakeupCommand 补课预约命令实现（内存版）。
type MakeupCommand struct {
	data *memory.Data
}

// 编译期断言：内存实现必须满足命令接口。
var _ repo.StudentMakeupCommand = (*MakeupCommand)(nil)

// NewMakeupCommand 创建补课预约命令实现
//
// 返回接口类型，这样实现不完整时会在编译期暴露。
func NewMakeupCommand(d *memory.Data) repo.StudentMakeupCommand {
	return &MakeupCommand{data: d}
}

// Save 保存补课预约（新增或更新）
func (c *MakeupCommand) Save(m *makeup.StudentMakeup) error {
	if m == nil {
		return repo.ErrMakeupRequired
	}
	c.data.SaveMakeup(m)
	return nil
}

// Delete 删除补课预约
func (c *MakeupCommand) Delete(id int64) error {
	if !c.data.DeleteMakeup(id) {
		return fmt.Errorf("%w: %d", repo.ErrMakeupNotFound, id)
	}
	return nil
}
