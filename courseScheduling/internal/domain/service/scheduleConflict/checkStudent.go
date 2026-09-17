package scheduleConflict

import "context"

// CheckStudent 检测「该学员在学课程的排期」与目标课程的排期是否撞时间。
//
// 占用集 = 他所有「在学」课程的全部排期；已退课 / 已结业的记录不计入
// （由 GetStudentSlots 负责过滤）。
//
// 重复选课也会被这里拦下：目标课程本身就在学员的在学课程里，两个集合里的
// 那份槽位是同一批，逐槽自比必然重叠。
func (s *Service) CheckStudent(ctx context.Context, studentID int64, courseID string) (bool, error) {
	target, err := s.SlotCmd.GetCourseSlots(ctx, courseID)
	if err != nil {
		return false, err
	}

	student, err := s.SlotCmd.GetStudentSlots(ctx, studentID)
	if err != nil {
		return false, err
	}

	return target.ConflictsWith(student), nil
}
