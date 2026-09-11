package query

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
)

// QualificationQuery 授课资质查询接口
//
// 资质绑定的是「课程类型」而不是具体课程，所以这里按课程类型维度查询。
type QualificationQuery interface {
	// Page 分页查询授课资质列表
	Page(ctx context.Context, page int, pageSize int) ([]*qualification.Qualification, error)
	// CourseTypesByTeacherID 根据讲师 ID 获取其有资质的课程类型列表
	CourseTypesByTeacherID(ctx context.Context, teacherID int64) ([]*courseType.CourseType, error)
	// TeachersByCourseTypeID 根据课程类型 ID 获取有资质的讲师列表
	TeachersByCourseTypeID(ctx context.Context, courseTypeID string) ([]*teacher.Teacher, error)
	// IsTeacherQualifiedForCourse 判断讲师是否有教指定课程的资质
	//
	// 链路：课程 → 所属课程类型 → 该讲师是否有该类型的资质。
	IsTeacherQualifiedForCourse(ctx context.Context, teacherID int64, courseID string) (bool, error)
}
