package qualificationCheck

import (
	"context"
	"fmt"

	"cqrs/internal/core/course_scheduling/domain/repo/command"
)

// CheckTeacherQualifiedForSlot 判断该讲师有没有资格教这一个槽位。
//
// 与 CheckTeacherQualified 的差别只在入参：那个收课程 ID，这个收槽位 ID
// —— 单槽位场景调用方手上通常只有槽位 ID。槽位 -> 课程这一段由本方法自己解析。
//
// 槽位不存在时报 not found —— 别把「传错槽位」当成「没资质」。
func (s *Service) CheckTeacherQualifiedForSlot(ctx context.Context, teacherID int64, slotID string) (bool, error) {
	slots, err := s.slots.GetSlots([]string{slotID})
	if err != nil {
		return false, err
	}
	if len(slots) == 0 {
		return false, fmt.Errorf("%w: %s", command.ErrCourseSlotNotFound, slotID)
	}

	return s.CheckTeacherQualified(ctx, teacherID, slots[0].CourseID())
}
