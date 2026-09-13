package query

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
)

// CourseQuery 课程查询接口
type CourseQuery interface {
	// PageCourses 分页查询课程列表
	PageCourses(ctx context.Context, page int, pageSize int) ([]*course.Course, error)
	// AvailableForStudent 查找与学员当前选课不冲突的课程
	AvailableForStudent(ctx context.Context, studentID int64) ([]*course.Course, error)
	// AvailableForTeacher 查找与讲师现有排课不冲突的课程
	AvailableForTeacher(ctx context.Context, teacherID int64) ([]*course.Course, error)
	// AvailableForClassroom 查找与教室现有排课不冲突的课程
	AvailableForClassroom(ctx context.Context, classroomID string) ([]*course.Course, error)
}
