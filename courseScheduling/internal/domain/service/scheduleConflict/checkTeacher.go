package scheduleConflict

import "context"

// CheckTeacher 检测「该讲师现有的课」与这批槽位是否撞时间。
//
// 讲师不存在时报 not found —— 调用方传错了 ID 应该报出来，
// 而不是当成「没冲突」悄悄放行。
func (s *Service) CheckTeacher(ctx context.Context, teacherID int64, slotIDs []string) (bool, error) {
	if _, err := s.teachers.MustGetTeacher(ctx, teacherID); err != nil {
		return false, err
	}

	existing, err := s.SlotCmd.GetTeacherSlots(ctx, teacherID)
	if err != nil {
		return false, err
	}

	slots, err := s.SlotCmd.GetSlots(ctx, slotIDs)
	if err != nil {
		return false, err
	}

	return existing.ConflictsWith(slots), nil
}
