package courseSlot

import (
	"context"

	command "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// AssignCourse 给指定课表槽位设置课程命令
type AssignCourse struct {
	SlotIDs  []string
	CourseID string
}

// AssignCourse 给指定课表槽位设置课程
//
// 仓库（命令适配器）会把 checkCourse 需要的数据装进 CourseAssignContext
// 一并传进来，所以这里直接把 checkCourse 当回调交给它。
func (h *Handler) AssignCourse(ctx context.Context, cmd AssignCourse) error {
	return h.SlotCmd.AssignCourse(ctx, cmd.SlotIDs, cmd.CourseID, checkCourse)
}

// checkCourse 判断把课程配到目标槽位是否存在冲突。
//
// 返回 true 表示有冲突（拒绝本次配置），false 表示可以配。
// 规则 = 与该课程现有槽位不撞时间。
func checkCourse(ctx context.Context, ac command.CourseAssignContext) (bool, error) {
	return ac.TargetSlots.ConflictsWith(ac.CourseSlots), nil
}
