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
// 把本次换课交给调用方判定（目标讲师/教室在目标时段是否已被占用），
// 通过后才写入；传 nil 表示不做检查。
func (c *CourseSlotChangeCommand) Change(ctx context.Context, csc *courseSlotChange.CourseSlotChange, checkConflictFn func(ctx context.Context, csc *courseSlotChange.CourseSlotChange) (bool, error)) error {
	if csc == nil {
		return repo.ErrSlotChangeRequired
	}

	if checkConflictFn != nil {
		conflict, err := checkConflictFn(ctx, csc)
		if err != nil {
			return err
		}
		if conflict {
			return fmt.Errorf("%w: course %s", repo.ErrSlotChangeConflict, csc.CourseID())
		}
	}
	return c.Save(csc)
}
