package command

import (
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
)

type AssignTeacherRepo interface {
	GetTeacherSlots(teacherID string) ([]courseSlot.CourseSlot, error)

	GetQualifications(teacherID string) ([]qualification.Qualification, error)

	GetCourseType(courseSlotsID string) (courseType.CourseType, error)
}
