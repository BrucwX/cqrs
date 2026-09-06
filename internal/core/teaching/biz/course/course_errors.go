package course

import "github.com/go-kratos/kratos/v3/errors"

var (
	// ErrCourseNotFound is returned when a course does not exist.
	ErrCourseNotFound = errors.NotFound("COURSE_NOT_FOUND", "course not found")
	// ErrCourseInvalidArgument is returned when a course request is invalid.
	ErrCourseInvalidArgument = errors.BadRequest("COURSE_INVALID_ARGUMENT", "invalid course argument")
	// ErrCourseTeacherNotFound is returned when the teacher assigned to a course does not exist.
	ErrCourseTeacherNotFound = errors.BadRequest("COURSE_TEACHER_NOT_FOUND", "teacher not found for course")
	// ErrCourseClassroomNotFound is returned when the classroom assigned to a course does not exist.
	ErrCourseClassroomNotFound = errors.BadRequest("COURSE_CLASSROOM_NOT_FOUND", "classroom not found for course")
	// ErrCourseInvalidSchedule is returned when the course schedule is invalid.
	ErrCourseInvalidSchedule = errors.BadRequest("COURSE_INVALID_SCHEDULE", "invalid course schedule")
	// ErrCourseFull is returned when the course has reached its maximum student capacity.
	ErrCourseFull = errors.BadRequest("COURSE_FULL", "course is full")
	// ErrCourseStudentAlreadyEnrolled is returned when a student is already enrolled in the course.
	ErrCourseStudentAlreadyEnrolled = errors.BadRequest("COURSE_STUDENT_ALREADY_ENROLLED", "student already enrolled")
	// ErrCourseStudentNotEnrolled is returned when a student is not enrolled in the course.
	ErrCourseStudentNotEnrolled = errors.BadRequest("COURSE_STUDENT_NOT_ENROLLED", "student not enrolled")
	// ErrCourseNoSessionsLeft is returned when all sessions have been consumed.
	ErrCourseNoSessionsLeft = errors.BadRequest("COURSE_NO_SESSIONS_LEFT", "no sessions left")
	// ErrCourseScheduleConflict is returned when the course schedule conflicts with another course.
	ErrCourseScheduleConflict = errors.BadRequest("COURSE_SCHEDULE_CONFLICT", "schedule conflicts with another course")
	// ErrCourseClassroomCapacityExceeded is returned when max students exceeds classroom capacity.
	ErrCourseClassroomCapacityExceeded = errors.BadRequest("COURSE_CLASSROOM_CAPACITY_EXCEEDED", "max students exceeds classroom capacity")
)
