package implement

import (
	"context"
	"fmt"
	"time"

	"cqrs/internal/core/course_scheduling/adapters/command/memory"
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
func (c *MakeupCommand) Save(ctx context.Context, m *makeup.StudentMakeup) error {
	if m == nil {
		return repo.ErrMakeupRequired
	}
	c.data.SaveMakeup(m)
	return nil
}

// Delete 删除补课预约
func (c *MakeupCommand) Delete(ctx context.Context, id int64) error {
	if !c.data.DeleteMakeup(id) {
		return fmt.Errorf("%w: %d", repo.ErrMakeupNotFound, id)
	}
	return nil
}

// GetMakeupsForTarget 取补到同一节课上的全部补课预约。
//
// 按槽位 + 日期两把钥匙匹配，状态不过滤。
func (c *MakeupCommand) GetMakeupsForTarget(ctx context.Context, targetSlotID string, targetDate time.Time) ([]makeup.StudentMakeup, error) {
	out := make([]makeup.StudentMakeup, 0)
	for _, item := range c.data.Makeups() {
		if item.TargetSlotID() != targetSlotID || !item.TargetDate().Equal(targetDate) {
			continue
		}
		out = append(out, *item)
	}
	return out, nil
}
