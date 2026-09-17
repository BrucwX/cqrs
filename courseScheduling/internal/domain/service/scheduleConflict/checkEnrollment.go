package scheduleConflict

import (
	"context"
	"time"
)

// CheckEnrollment 判断该学员能不能选这门课。
//
// 规则 = 课程还开着选课窗口 且 没满员 且 目标课程的排期跟学员在学课程的排期不撞时间。
//
// 课程不存在时报 not found —— 调用方传错了 ID 应该报出来，
// 而不是当成「不能选」悄悄拦下。
func (s *Service) CheckEnrollment(ctx context.Context, studentID int64, courseID string) (bool, error) {
	crs, err := s.courses.MustGetCourse(ctx, courseID)
	if err != nil {
		return false, err
	}

	// 不在窗口内 / 已满 -> 不能选
	if err := crs.CanEnroll(time.Now()); err != nil {
		return true, nil
	}

	// 重复选课也会被下面这条拦下：目标课程本身就在他的在学课程里，
	// 两个集合里那份槽位是同一批，逐槽自比必然重叠。
	return s.CheckStudent(ctx, studentID, courseID)
}
