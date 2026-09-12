package courseSlotChange

import (
	"context"
	"fmt"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlotChange"
)

// ChangeCourseSlot 临时换课命令
type ChangeCourseSlot struct {
	CourseID    string
	ApplicantID int64
	ChangeType  courseSlotChange.ChangeType
	Original    courseSlotChange.OriginalPlan
	Target      courseSlotChange.TargetPlan
	Reason      string
}

// ChangeCourseSlot 临时换课（调课 / 代课 / 换教室）
//
// 登记即生效，没有审批环节；参数校验由聚合根构造函数负责
// （课程/申请人/事由必填，目标时间不能早于现在，目标起止时间与讲师、教室必填）。
// 顺序：先判、过了才写。判定交给领域服务 ScheduleConflict.CheckSlotChange，
// 写回走 ChangeCmd.Change。
func (h *Handler) ChangeCourseSlot(ctx context.Context, cmd ChangeCourseSlot) (change *courseSlotChange.CourseSlotChange, err error) {
	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { err = h.tx.End(ctx, err) }()

	change, err = courseSlotChange.NewCourseSlotChange(
		cmd.CourseID,
		cmd.ApplicantID,
		cmd.ChangeType,
		cmd.Original,
		cmd.Target,
		cmd.Reason,
		time.Now(),
	)
	if err != nil {
		return nil, err
	}

	// 判定：目标讲师 / 教室在目标时段没被占用，也没有别的换课占在那里
	conflict, err := h.conflict.CheckSlotChange(ctx, change)
	if err != nil {
		return nil, err
	}
	if conflict {
		return nil, fmt.Errorf("%w: course %s", courseSlotChange.ErrSlotChangeConflict, cmd.CourseID)
	}

	if err := h.ChangeCmd.Change(ctx, change); err != nil {
		return nil, err
	}

	return change, nil
}
