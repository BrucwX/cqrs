package implement

import (
	"context"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/command/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// CourseSlotChangeCommand 课表变更命令实现（内存版）。
type CourseSlotChangeCommand struct {
	data *memory.Data
}

// 编译期断言：内存实现必须满足命令接口。
var _ repo.CourseSlotChangeCommand = (*CourseSlotChangeCommand)(nil)

// NewCourseSlotChangeCommand 创建课表变更命令实现
//
// 返回接口类型，这样实现不完整时会在编译期暴露。
func NewCourseSlotChangeCommand(d *memory.Data) repo.CourseSlotChangeCommand {
	return &CourseSlotChangeCommand{data: d}
}

// Save 保存课表变更（新增或更新）
func (c *CourseSlotChangeCommand) Save(csc *courseSlotChange.CourseSlotChange) error {
	if csc == nil {
		return repo.ErrSlotChangeRequired
	}
	c.data.SaveCourseSlotChange(csc)
	return nil
}

// Delete 删除课表变更
func (c *CourseSlotChangeCommand) Delete(id int64) error {
	if !c.data.DeleteCourseSlotChange(id) {
		return fmt.Errorf("%w: %d", repo.ErrSlotChangeNotFound, id)
	}
	return nil
}

// Change 登记一次临时换课
//
// 只管写：冲突判定（目标讲师 / 教室在目标时段是否已被占用）由调用方在调过来之前做完。
func (c *CourseSlotChangeCommand) Change(ctx context.Context, csc *courseSlotChange.CourseSlotChange) error {
	if csc == nil {
		return repo.ErrSlotChangeRequired
	}

	return c.Save(csc)
}

// GetOtherSlotChanges 取除该变更单以外的全部换课记录。
func (c *CourseSlotChangeCommand) GetOtherSlotChanges(id int64) ([]*courseSlotChange.CourseSlotChange, error) {
	all := c.data.CourseSlotChanges()

	out := make([]*courseSlotChange.CourseSlotChange, 0, len(all))
	for _, item := range all {
		if item.ID() != id {
			out = append(out, item)
		}
	}
	return out, nil
}
