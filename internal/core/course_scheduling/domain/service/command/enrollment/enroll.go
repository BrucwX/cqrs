package enrollment

import (
	"context"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
)

// StudentEnroll 学生选课命令
type StudentEnroll struct {
	StudentID int64
	CourseID  string
}

// StudentEnroll 学生选课
//
// 仓库（命令适配器）会把「待写入的注册记录 + 该课程聚合」传进来，判定走 h.checkEnrollment。
func (h *Handler) StudentEnroll(ctx context.Context, cmd StudentEnroll) (*enrollment.CourseEnrollment, error) {
	// 创建选课记录
	enroll, err := enrollment.NewCourseEnrollment(cmd.StudentID, cmd.CourseID)
	if err != nil {
		return nil, err
	}

	// 保存（带选课准入检查）
	if err := h.EnrollmentCmd.Enroll(ctx, enroll, h.checkEnrollment); err != nil {
		return nil, err
	}

	return enroll, nil
}

// checkEnrollment 判断这次选课是否存在冲突。
//
// 返回 true 表示有冲突（拒绝本次选课），false 表示可以选。
// 规则 = 在选课窗口内 且 课程未满 且 与学员现有在学课程不撞时间。
//
// 重复选课也会被时间那一条拦下：学员现有在学课程里包含目标课程本身，
// 它的槽位跟目标槽位同一份，逐槽自比必然重叠。
//
// 仓库只把「待写入的注册记录 + 该课程聚合」传进来，排期按需自己取。
func (h *Handler) checkEnrollment(ctx context.Context, e *enrollment.CourseEnrollment, c course.Course) (bool, error) {
	// 1) 选课窗口 + 容量：交给课程聚合自己的准入规则
	if err := c.CanEnroll(time.Now()); err != nil {
		return true, nil // 不在窗口内 / 已满 -> 冲突
	}

	// 2) 时间：与学员现有在学课程排期重叠 -> 冲突
	targetSlots, err := h.enroll.GetCourseSlots(e.CourseID())
	if err != nil {
		return false, err
	}

	studentSlots, err := h.enroll.GetStudentSlots(e.StudentID())
	if err != nil {
		return false, err
	}

	return targetSlots.ConflictsWith(studentSlots), nil
}
