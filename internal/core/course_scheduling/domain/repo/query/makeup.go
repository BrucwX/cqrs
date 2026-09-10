package query

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/makeup"
	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
)

// StudentMakeupQuery 补课申请查询接口
type StudentMakeupQuery interface {
	// Page 分页查询补课申请列表
	Page(ctx context.Context, page int, pageSize int) ([]*makeup.StudentMakeup, error)
	// CoursesByStudentID 根据学员 ID 获取已补课的课程列表
	CoursesByStudentID(ctx context.Context, studentID int64) ([]*course.Course, error)
	// StudentsByCourseID 根据课程 ID 获取已补课的学员列表
	StudentsByCourseID(ctx context.Context, courseID string) ([]*student.Student, error)
}
