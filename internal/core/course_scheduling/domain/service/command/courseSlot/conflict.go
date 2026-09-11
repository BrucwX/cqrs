package courseSlot

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
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
func checkTeacherQualification(ctx context.Context, cty courseType.CourseType, t_q []qualification.Qualification) (bool, error) {
	// 检查老师 ID是否匹配
	for _, q := range t_q {
		if q.CourseTypeID() == cty.ID() {
			return false, nil // 老师有资质，没有冲突
		}
	}

	return true, nil // 老师没有资质，有冲突
}

// checkClassroomCapacity 检查教室容量是否能容纳这门课
// 返回 true 表示教室容量不足（有冲突），false 表示教室容量足够（没有冲突）
func checkClassroomCapacity(ctx context.Context, c course.Course, cl classroom.Classroom) (bool, error) {
	// 课程最大人数 <= 教室容量，教室足够大，没有冲突
	if c.Capacity().Max() < cl.Capacity() {
		return false, nil
	}
	// 教室容量不足以容纳课程，有冲突
	return true, nil
}
