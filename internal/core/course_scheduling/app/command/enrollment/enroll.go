package enrollment

import (
	"context"
	"fmt"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
)

// StudentEnroll 学生选课命令
type StudentEnroll struct {
	StudentID int64
	CourseID  string
}

// StudentEnroll 学生选课
//
// 不新建注册记录：记录是报名登记（缴费）时落下的「未选课」资格，
// 这里只把它推到在读。顺序：① 找未选课记录（没有 = 没资格）→ ② 判冲突 → ③ 改状态写回。
func (h *Handler) StudentEnroll(ctx context.Context, cmd StudentEnroll) (enroll *enrollment.CourseEnrollment, err error) {
	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { err = h.tx.End(ctx, err) }()

	// 1) 必须已经有一条「未选课」的注册记录，没有就是这门课还没缴费
	//    —— 未付费错误由仓库直接报出来（ErrEnrollmentNotPaid）。
	pending, err := h.EnrollmentCmd.GetUnSelectEnrollBySC(ctx, cmd.StudentID, cmd.CourseID)
	if err != nil {
		return nil, err
	}

	// 2) 判定：选课窗口 + 容量 + 时间冲突。
	//    课程不存在也在这里报 not found —— 服务里那一下 GetCourse 顺带校了。
	rejected, err := h.conflict.CheckEnrollment(ctx, cmd.StudentID, cmd.CourseID)
	if err != nil {
		return nil, err
	}
	if rejected {
		return nil, fmt.Errorf("%w: student %d course %s", enrollment.ErrEnrollmentConflict, cmd.StudentID, cmd.CourseID)
	}

	// 3) 未选课 -> 在读，写回同一条记录
	if err := pending.Enroll(time.Now()); err != nil {
		return nil, err
	}
	if err := h.EnrollmentCmd.Enroll(ctx, &pending); err != nil {
		return nil, err
	}

	return &pending, nil
}
