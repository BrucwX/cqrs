package command

import (
	"errors"

	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
)

// ErrCourseNotFinished 讲师还没以学员身份修完该课程类型下的课程（没结业，或缺过课）。
var ErrCourseNotFinished = errors.New("teacher has not finished a course of this type")

// QualifyRepo 是授予授课资质做准入检查所需的数据来源。
//
// 判定（「修完该类型下的某门课程」）由命令服务自己实现，它需要
// 「该讲师聚合 / 该课程类型下的课程 / 该讲师作为学员的报名与缺勤记录」；
// 仓库只负责按需提供这些数据，不参与规则判定。
type QualifyRepo interface {
	// GetTeacher 取该讲师聚合（要用它的学员 ID 去查报名/缺勤）
	GetTeacher(teacherID int64) (teacher.Teacher, error)
	// GetCourses 取某课程类型下的全部课程
	GetCourses(courseTypeID string) ([]course.Course, error)
	// GetEnrollments 取该学员（讲师以学员身份参训）的全部报名记录
	GetEnrollments(studentID int64) ([]enrollment.CourseEnrollment, error)
	// GetAbsences 取该学员的全部缺勤记录
	GetAbsences(studentID int64) ([]absence.AbsenceRecord, error)
}
