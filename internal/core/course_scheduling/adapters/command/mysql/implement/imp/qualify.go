package imp

import (
	"cqrs/internal/core/course_scheduling/adapters/command/mysql"
	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

type QualifyImp struct {
	data *mysql.Data
}

var _ repo.QualifyRepo = (*QualifyImp)(nil)

func NewQualifyImp(d *mysql.Data) repo.QualifyRepo {
	return &QualifyImp{data: d}
}

func (c *QualifyImp) GetTeacher(teacherID int64) (teacher.Teacher, error) {
	panic("implement me")
}

func (c *QualifyImp) GetCourses(courseTypeID string) ([]course.Course, error) {
	panic("implement me")
}

func (c *QualifyImp) GetEnrollments(studentID int64) ([]enrollment.CourseEnrollment, error) {
	panic("implement me")
}

func (c *QualifyImp) GetAbsences(studentID int64) ([]absence.AbsenceRecord, error) {
	panic("implement me")
}
