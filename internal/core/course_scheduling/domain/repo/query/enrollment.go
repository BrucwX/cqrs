package query

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
)

// CourseEnrollmentQuery 课程注册查询接口
type CourseEnrollmentQuery interface {
	// PageEnrollments 分页查询课程注册列表
	PageEnrollments(ctx context.Context, page int, pageSize int) ([]*enrollment.CourseEnrollment, error)
	// EnrolledCoursesByStudentID 根据学员 ID 获取所注册的课程列表
	EnrolledCoursesByStudentID(ctx context.Context, studentID int64) ([]*course.Course, error)
	// EnrolledStudentsByCourseID 根据课程 ID 获取所有注册的学员列表
	EnrolledStudentsByCourseID(ctx context.Context, courseID string) ([]*student.Student, error)
}
