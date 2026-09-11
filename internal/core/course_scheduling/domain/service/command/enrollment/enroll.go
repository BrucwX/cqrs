package enrollment

import (
	"context"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
	command "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// StudentEnroll 学生选课命令
type StudentEnroll struct {
	StudentID int64
	CourseID  string
}

// StudentEnroll 学生选课
//
// 冲突判定走 checkEnrollment：仓库（命令适配器）会把「目标课程 / 目标课程排期 /
// 学员现有排期」装进 EnrollmentContext 一并传进来。
func (h *Handler) StudentEnroll(ctx context.Context, cmd StudentEnroll) (*enrollment.CourseEnrollment, error) {
	// 创建选课记录
	enroll, err := enrollment.NewCourseEnrollment(cmd.StudentID, cmd.CourseID)
	if err != nil {
		return nil, err
	}

	// 保存（带选课准入检查）
	if err := h.EnrollmentCmd.Enroll(ctx, enroll, checkEnrollment); err != nil {
		return nil, err
	}

	return enroll, nil
}

// checkEnrollment 判断这次选课是否存在冲突。
//
// 返回 true 表示有冲突（拒绝本次选课），false 表示可以选。
// 规则 = 在选课窗口内 且 课程未满 且 与学员现有在学课程不撞时间。
//
// 重复选课也会被第 3 条拦下：学员现有在学课程里包含目标课程本身，
// 它的槽位跟目标槽位同一份，逐槽自比必然重叠。
func checkEnrollment(ctx context.Context, ec command.EnrollmentContext) (bool, error) {
	// 1) 选课窗口 + 容量：交给课程聚合自己的准入规则
	if err := ec.Course.CanEnroll(time.Now()); err != nil {
		return true, nil // 不在窗口内 / 已满 -> 冲突
	}

	// 2) 时间：与学员现有在学课程排期重叠 -> 冲突
	return ec.TargetSlots.ConflictsWith(ec.StudentSlots), nil
}
