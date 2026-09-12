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
func (c *CourseSlotChangeCommand) Save(ctx context.Context, csc *courseSlotChange.CourseSlotChange) error {
	if csc == nil {
		return courseSlotChange.ErrSlotChangeRequired
	}
	c.data.SaveCourseSlotChange(csc)
	return nil
}

// Delete 删除课表变更
func (c *CourseSlotChangeCommand) Delete(ctx context.Context, id int64) error {
	if !c.data.DeleteCourseSlotChange(id) {
		return fmt.Errorf("%w: %d", courseSlotChange.ErrSlotChangeNotFound, id)
	}
	return nil
}

// MustGet 取课表变更聚合；不存在时报 ErrSlotChangeNotFound。
func (c *CourseSlotChangeCommand) MustGet(ctx context.Context, id int64) (courseSlotChange.CourseSlotChange, error) {
	for _, item := range c.data.CourseSlotChanges() {
		if item.ID() == id {
			return *item, nil
		}
	}
	return courseSlotChange.CourseSlotChange{}, fmt.Errorf("%w: %d", courseSlotChange.ErrSlotChangeNotFound, id)
}

// Change 登记一次临时换课
//
// 只管写：冲突判定（目标讲师 / 教室在目标时段是否已被占用）由调用方在调过来之前做完。
func (c *CourseSlotChangeCommand) Change(ctx context.Context, csc *courseSlotChange.CourseSlotChange) error {
	if csc == nil {
		return courseSlotChange.ErrSlotChangeRequired
	}

	return c.Save(ctx, csc)
}

// GetOtherSlotChanges 取除该变更单以外的全部换课记录。
func (c *CourseSlotChangeCommand) GetOtherSlotChanges(ctx context.Context, id int64) ([]*courseSlotChange.CourseSlotChange, error) {
	all := c.data.CourseSlotChanges()

	out := make([]*courseSlotChange.CourseSlotChange, 0, len(all))
	for _, item := range all {
		if item.ID() != id {
			out = append(out, item)
		}
	}
	return out, nil
}
