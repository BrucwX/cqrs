package biz

import (
	"context"

	"cqrs/internal/core/teaching/biz/course"
)

// ScheduleValidator validates schedule conflicts across aggregates.
type ScheduleValidator struct {
	finder *CourseFinder
}

// NewScheduleValidator creates a new ScheduleValidator.
func NewScheduleValidator(finder *CourseFinder) *ScheduleValidator {
	return &ScheduleValidator{
		finder: finder,
	}
}

// ValidateTeacherScheduleConflict checks if the given weekly slots conflict
// with any existing courses of the specified teacher within the semester.
func (v *ScheduleValidator) ValidateTeacherScheduleConflict(
	ctx context.Context,
	teacherID string,
	slots []course.ScheduleTime,
	semester course.Semester,
) error {
	courses, err := v.finder.ListCoursesByTeacher(ctx, teacherID)
	if err != nil {
		return err
	}
	for _, c := range courses {
		if hasScheduleConflict(slots, semester, c) {
			return course.ErrCourseScheduleConflict
		}
	}
	return nil
}

// ValidateClassroomScheduleConflict checks if the given weekly slots conflict
// with any existing courses of the specified classroom within the semester.
func (v *ScheduleValidator) ValidateClassroomScheduleConflict(
	ctx context.Context,
	classroomID string,
	slots []course.ScheduleTime,
	semester course.Semester,
) error {
	courses, err := v.finder.ListCoursesByClassroom(ctx, classroomID)
	if err != nil {
		return err
	}
	for _, c := range courses {
		if hasScheduleConflict(slots, semester, c) {
			return course.ErrCourseScheduleConflict
		}
	}
	return nil
}

// hasScheduleConflict checks if the given slots conflict with an existing course.
func hasScheduleConflict(
	slots []course.ScheduleTime,
	semester course.Semester,
	existing *course.Course,
) bool {
	// 学期日期范围不重叠则无冲突
	if !semester.Overlaps(&existing.Semester) {
		return false
	}
	// 检查每周时间段是否有重叠
	for _, newSlot := range slots {
		for _, existingSlot := range existing.WeeklySlots {
			if newSlot.Overlaps(&existingSlot) {
				return true
			}
		}
	}
	return false
}
