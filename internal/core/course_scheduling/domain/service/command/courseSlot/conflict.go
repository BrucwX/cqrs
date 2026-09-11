package courseSlot

import (
	"context"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
)

// ConflictChecker 冲突检查函数类型
type ConflictChecker func(ctx context.Context, targetSlots []courseSlot.CourseSlot, existingSlots []courseSlot.CourseSlot) (bool, error)

// checkScheduleConflict 检查目标课表与现有课表是否有冲突
// 通用的冲突检查逻辑：检查同一天的时间段是否有重叠
func checkScheduleConflict(ctx context.Context, targetSlots []courseSlot.CourseSlot, existingSlots []courseSlot.CourseSlot) (bool, error) {
	// 遍历目标课表
	for _, targetSlot := range targetSlots {
		// 遍历现有课表
		for _, existingSlot := range existingSlots {
			// 检查是否同一天
			if targetSlot.Weekday() != existingSlot.Weekday() {
				continue
			}
			// 检查时间段是否有重叠
			if targetSlot.TimeRange().Overlaps(existingSlot.TimeRange()) {
				return true, nil // 有冲突
			}
		}
	}
	return false, nil // 没有冲突
}

// checkTeacherQualification 检查老师是否有能力上这门课
// 返回 true 表示老师没有资质（有冲突），false 表示老师有资质（没有冲突）
func checkTeacherQualification(ctx context.Context, course course.Course, teacher teacher.Teacher, qualification qualification.Qualification) (bool, error) {
	// 检查老师 ID是否匹配
	if qualification.TeacherID() != teacher.ID() {
		return true, nil // 老师 ID不匹配，有冲突
	}

	// 检查课程 ID是否匹配
	if qualification.CourseID() != course.ID() {
		return true, nil // 课程 ID不匹配，有冲突
	}

	// 检查资质是否有效
	if err := qualification.IsEligible(time.Now()); err != nil {
		return true, nil // 资质无效（过期或被吊销），有冲突
	}

	return false, nil // 老师有资质，没有冲突
}
