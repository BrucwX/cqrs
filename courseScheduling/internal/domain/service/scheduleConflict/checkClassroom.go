package scheduleConflict

import "context"

// CheckClassroom 检测「该教室现有的课」与这批槽位是否撞时间。
func (s *Service) CheckClassroom(ctx context.Context, classroomID string, slotIDs []string) (bool, error) {
	existing, err := s.SlotCmd.GetClassroomSlots(ctx, classroomID)
	if err != nil {
		return false, err
	}

	slots, err := s.SlotCmd.GetSlots(ctx, slotIDs)
	if err != nil {
		return false, err
	}

	return existing.ConflictsWith(slots), nil
}
