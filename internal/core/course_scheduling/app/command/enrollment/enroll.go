package enrollment

import (
	"context"
	"fmt"

	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
)

// StudentEnroll 学生选课命令
type StudentEnroll struct {
	StudentID int64
	CourseID  string
}

// StudentEnroll 学生选课
//
// 顺序：先判、过了才写。判定交给领域服务 ScheduleConflict.CheckEnrollment，
// 写回走 EnrollmentCmd.Enroll。
func (h *Handler) StudentEnroll(ctx context.Context, cmd StudentEnroll) (enroll *enrollment.CourseEnrollment, err error) {
	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { err = h.tx.End(ctx, err) }()

	// 1) 创建选课记录（聚合构造函数负责填参校验）
	enroll, err = enrollment.NewCourseEnrollment(cmd.StudentID, cmd.CourseID)
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

	// 3) 写回
	if err := h.EnrollmentCmd.Enroll(ctx, enroll); err != nil {
		return nil, err
	}

	return enroll, nil
}
