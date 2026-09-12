package scheduleConflict

import "context"

// CheckCourse 检测「该课程现有的课」与这批槽位是否撞时间。
//
// 课程不存在时报 not found —— 调用方传错了 ID 应该报出来，
// 而不是当成「没冲突」悄悄放行（传错 ID 的槽位被写进去就指向空课程了）。
//
// 目标槽位本身已经属于这门课的话同样会判成冲突 —— 那说明重复配置了，
// 应该报出来让调用方处理。
func (s *Service) CheckCourse(ctx context.Context, courseID string, slotIDs []string) (bool, error) {
	if _, err := s.courses.GetCourse(courseID); err != nil {
		return false, err
	}

	existing, err := s.SlotCmd.GetCourseSlots(courseID)
	if err != nil {
		return false, err
	}

	slots, err := s.SlotCmd.GetSlots(slotIDs)
	if err != nil {
		return false, err
	}

	return existing.ConflictsWith(slots), nil
}
