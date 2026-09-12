package classroomCapacity

import "context"

// Check 判断这批槽位能不能都放进目标教室。
//
// 教室由调用方给（挑教室的场景），逐槽判：任一槽位所属课程的人数超过教室容量
// 就拦下，不管其余槽位过不过。
//
// 教室不存在时报 not found —— 调用方传错了 ID 应该报出来，
// 而不是当成「装不下」或者「装得下」悄悄放行。
func (s *Service) Check(ctx context.Context, classroomID string, slotIDs []string) (bool, error) {
	cr, err := s.classrooms.MustGet(ctx, classroomID)
	if err != nil {
		return false, err
	}

	slots, err := s.slots.GetSlots(ctx, slotIDs)
	if err != nil {
		return false, err
	}

	for _, slot := range slots {
		crs, err := s.courses.MustGet(ctx, slot.CourseID())
		if err != nil {
			return false, err
		}

		// 教室自带的方法，还会判教室状态与已分配座位
		if err := cr.CanAccommodate(crs.Capacity().Max()); err != nil {
			return true, nil
		}
	}

	return false, nil
}
