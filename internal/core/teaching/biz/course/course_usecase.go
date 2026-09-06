package course

import (
	"github.com/go-kratos/kratos-layout/internal/core/teaching/biz/teacher"
	"github.com/go-kratos/kratos-layout/internal/core/venue/biz/classroom"
)

// CourseUsecase is the course usecase.
type CourseUsecase struct {
	Repo          CourseRepo
	TeacherRepo   teacher.TeacherRepo
	ClassroomRepo classroom.ClassroomRepo
}

// NewCourseUsecase creates a new CourseUsecase.
func NewCourseUsecase(repo CourseRepo, teacherRepo teacher.TeacherRepo, classroomRepo classroom.ClassroomRepo) *CourseUsecase {
	return &CourseUsecase{
		Repo:          repo,
		TeacherRepo:   teacherRepo,
		ClassroomRepo: classroomRepo,
	}
}
