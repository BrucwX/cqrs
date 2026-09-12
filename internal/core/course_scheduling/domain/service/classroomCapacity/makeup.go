package classroomCapacity

import (
	"context"
	"fmt"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
	"cqrs/internal/core/course_scheduling/domain/repo/command"
)

// CheckMakeup 判断学员补到某节课上时，那节课的教室装不装得下。
//
// 要坐进去的人 = 该课程的人数上限（正式学员）
// + 已经约在同一节课上的其他补课学员（只算 StatusBooked）
// + 这次要新增的这一位。
//
// 「同一节课」按槽位 + 日期两把钥匙认：同一个周排槽位落在别的日期上是另一堂课，
// 不抢这间教室的座位。已取消的不会来；已核销的是过去那节，也不算。
//
// 目标槽位不存在时报 not found；槽位还没排教室时报 classroom not found ——
// 两种都是「传错 ID」或「这节排期还没排好」，不该当成「装不下」悄悄拦下。
func (s *Service) CheckMakeup(ctx context.Context, targetSlotID string, targetDate time.Time) (bool, error) {
	slots, err := s.slots.GetSlots(ctx, []string{targetSlotID})
	if err != nil {
		return false, err
	}
	if len(slots) == 0 {
		return false, fmt.Errorf("%w: %s", command.ErrCourseSlotNotFound, targetSlotID)
	}
	target := slots[0]

	crs, err := s.courses.MustGet(ctx, target.CourseID())
	if err != nil {
		return false, err
	}

	cr, err := s.classrooms.MustGet(ctx, target.ClassroomID())
	if err != nil {
		return false, err
	}

	// 这门课的正式学员 + 这位新来的
	seats := crs.Capacity().Max() + 1

	// 再加已经约在同一节课上的其他补课学员，只算还等着上课的
	others, err := s.makeups.GetMakeupsForTarget(ctx, targetSlotID, targetDate)
	if err != nil {
		return false, err
	}
	for _, other := range others {
		if other.Status() == makeup.StatusBooked {
			seats++
		}
	}

	// 教室自带的方法，还会判教室状态与已分配座位
	if err := cr.CanAccommodate(seats); err != nil {
		return true, nil
	}

	return false, nil
}
