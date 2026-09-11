package command

import (
	"context"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/memory"
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
// 先把 checkSlotChange 需要的上下文装好交给调用方判定（目标讲师/教室在目标时段
// 是否已被占用），通过后才写入；传 nil 表示不做检查。
func (c *CourseSlotChangeCommand) Change(ctx context.Context, csc *courseSlotChange.CourseSlotChange, checkConflictFn func(ctx context.Context, sc repo.SlotChangeContext) (bool, error)) error {
	if csc == nil {
		return repo.ErrSlotChangeRequired
	}

	if checkConflictFn != nil {
		conflict, err := checkConflictFn(ctx, c.slotChangeContext(csc))
		if err != nil {
			return err
		}
		if conflict {
			return fmt.Errorf("%w: course %s", repo.ErrSlotChangeConflict, csc.CourseID())
		}
	}

	return c.Save(csc)
}

// --- 内部实现 ---

// slotChangeContext 组装换课时冲突检查所需的上下文。
func (c *CourseSlotChangeCommand) slotChangeContext(csc *courseSlotChange.CourseSlotChange) repo.SlotChangeContext {
	target := csc.TargetPlan()

	ctx := repo.SlotChangeContext{
		Change:       *csc,
		OtherChanges: c.otherChanges(csc.ID()),
	}

	for _, cs := range c.data.CourseSlots() {
		// 本门课自己的排期就是被换掉的那节课，不算占用
		if cs.CourseID() == csc.CourseID() {
			continue
		}
		if cs.TeacherID() == target.TeacherID() {
			ctx.TeacherSlots = append(ctx.TeacherSlots, *cs)
		}
		if cs.ClassroomID() == target.ClassroomID() {
			ctx.ClassroomSlots = append(ctx.ClassroomSlots, *cs)
		}
	}

	return ctx
}

// otherChanges 取除自己以外的全部换课记录。
func (c *CourseSlotChangeCommand) otherChanges(id int64) []*courseSlotChange.CourseSlotChange {
	all := c.data.CourseSlotChanges()
	out := make([]*courseSlotChange.CourseSlotChange, 0, len(all))
	for _, item := range all {
		if item.ID() != id {
			out = append(out, item)
		}
	}
	return out
}
