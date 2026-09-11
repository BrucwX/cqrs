package courseSlot

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
)

// AssignCourseInput 给指定课表槽位设置课程命令
type AssignCourseInput struct {
	SlotIDs  []string
	CourseID string
}

// AssignCourse 给指定课表槽位设置课程
//
// 仓库（命令适配器）会把「本次要配的槽位 + 该课程聚合」传进来，判定走 h.checkCourse。
func (h *Handler) AssignCourse(ctx context.Context, cmd AssignCourseInput) error {
	return h.SlotCmd.AssignCourse(ctx, cmd.SlotIDs, cmd.CourseID, h.checkCourse)
}

// checkCourse 判断把课程配到目标槽位是否存在冲突。
//
// 返回 true 表示有冲突（拒绝本次配置），false 表示可以配。
// 规则 = 与该课程现有排课不撞时间。
//
// 目标槽位本身已经属于这门课的话同样会判成冲突 —— 那说明重复配置了，
// 应该报出来让调用方处理。
//
// 仓库只把「本次要配的槽位 + 该课程聚合」传进来，课程排期按需自己取。
func (h *Handler) checkCourse(ctx context.Context, slots []courseSlot.CourseSlot, c course.Course) (bool, error) {
	courseSlots, err := h.assignCourse.GetCourseSlots(c.ID())
	if err != nil {
		return false, err
	}

	return courseSlots.ConflictsWith(courseSlot.CourseSlots(slots)), nil
}
