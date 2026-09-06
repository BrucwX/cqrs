package service

import (
	"cqrs/internal/core/teaching/biz/course"
)

// CourseService is a course service.
type CourseService struct {
	uc *course.CourseUsecase
}

// NewCourseService creates a new CourseService.
func NewCourseService(uc *course.CourseUsecase) *CourseService {
	return &CourseService{uc: uc}
}
