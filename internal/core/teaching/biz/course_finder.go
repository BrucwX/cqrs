package biz

import (
	"context"

	"cqrs/internal/core/teaching/biz/course"
)

// CourseFinder queries existing courses to build conflict-candidate sets for
// scheduling decisions (e.g. 课选教室). It reads through the course repo and
// filters by a single dimension: teacher or classroom.
type CourseFinder struct {
	courseRepo course.CourseRepo
}

// NewCourseFinder creates a new CourseFinder.
func NewCourseFinder(courseRepo course.CourseRepo) *CourseFinder {
	return &CourseFinder{
		courseRepo: courseRepo,
	}
}

// ListCoursesByTeacher returns all existing courses taught by the given teacher.
// It is the candidate set when checking whether a new course would conflict
// with the teacher's existing schedule.
func (f *CourseFinder) ListCoursesByTeacher(ctx context.Context, teacherID string) ([]*course.Course, error) {
	allCourses, err := f.courseRepo.ListCourses(ctx)
	if err != nil {
		return nil, err
	}
	matched := make([]*course.Course, 0)
	for _, c := range allCourses {
		if c.TeacherID == teacherID {
			matched = append(matched, c)
		}
	}
	return matched, nil
}

// ListCoursesByClassroom returns all existing courses held in the given classroom.
// It is the candidate set when checking whether a new course would conflict
// with the classroom's existing bookings.
func (f *CourseFinder) ListCoursesByClassroom(ctx context.Context, classroomID string) ([]*course.Course, error) {
	allCourses, err := f.courseRepo.ListCourses(ctx)
	if err != nil {
		return nil, err
	}
	matched := make([]*course.Course, 0)
	for _, c := range allCourses {
		if c.ClassroomID == classroomID {
			matched = append(matched, c)
		}
	}
	return matched, nil
}
