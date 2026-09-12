package implement

import (
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/command/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// QualifyRepo 是 repo.QualifyRepo 的内存实现。
//
// 它只做「按需取数据」：授证的准入判定在命令服务里（checkFinished），
// 这里不参与任何规则。
type QualifyRepo struct {
	data *memory.Data
}

// 编译期断言：内存实现必须满足接口。
var _ repo.QualifyRepo = (*QualifyRepo)(nil)

// NewQualifyRepo 创建授证准入检查的数据来源。
func NewQualifyRepo(d *memory.Data) repo.QualifyRepo {
	return &QualifyRepo{data: d}
}

// GetTeacher 取该讲师聚合，取不到报 not found。
func (r *QualifyRepo) GetTeacher(teacherID int64) (teacher.Teacher, error) {
	item, ok := r.data.TeacherByID(teacherID)
	if !ok {
		return teacher.Teacher{}, fmt.Errorf("%w: %d", repo.ErrTeacherNotFound, teacherID)
	}
	return *item, nil
}

// GetCourses 取某课程类型下的全部课程。
func (r *QualifyRepo) GetCourses(courseTypeID string) ([]course.Course, error) {
	out := make([]course.Course, 0)
	for _, item := range r.data.Courses() {
		if item.CourseTypeID() == courseTypeID {
			out = append(out, *item)
		}
	}
	return out, nil
}

// GetEnrollments 取该学员（讲师以学员身份参训）的全部报名记录。
func (r *QualifyRepo) GetEnrollments(studentID int64) ([]enrollment.CourseEnrollment, error) {
	out := make([]enrollment.CourseEnrollment, 0)
	for _, item := range r.data.Enrollments() {
		if item.StudentID() == studentID {
			out = append(out, *item)
		}
	}
	return out, nil
}

// GetAbsences 取该学员的全部缺勤记录。
func (r *QualifyRepo) GetAbsences(studentID int64) ([]absence.AbsenceRecord, error) {
	out := make([]absence.AbsenceRecord, 0)
	for _, item := range r.data.Absences() {
		if item.StudentID() == studentID {
			out = append(out, *item)
		}
	}
	return out, nil
}
