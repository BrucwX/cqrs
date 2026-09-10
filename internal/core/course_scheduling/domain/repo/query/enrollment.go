package query

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
)

// CourseEnrollmentQuery 课程注册查询接口
type CourseEnrollmentQuery interface {
	// Page 分页查询课程注册列表
	Page(ctx context.Context, page int, pageSize int) ([]*enrollment.CourseEnrollment, error)
	// CoursesByStudentID 根据学员 ID 获取所注册的课程列表
	CoursesByStudentID(ctx context.Context, studentID int64) ([]*course.Course, error)
	// StudentsByCourseID 根据课程 ID 获取所有注册的学员列表
	StudentsByCourseID(ctx context.Context, courseID string) ([]*student.Student, error)
}
