package service

import (
	"cqrs/internal/core/teaching/biz/student"
)

// StudentService is a student service.
type StudentService struct {
	uc *student.StudentUsecase
}

// NewStudentService creates a new StudentService.
func NewStudentService(uc *student.StudentUsecase) *StudentService {
	return &StudentService{uc: uc}
}
