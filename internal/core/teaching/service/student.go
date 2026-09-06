package service

import (
	"github.com/go-kratos/kratos-layout/internal/core/teaching/biz/student"
)

// StudentService is a student service.
type StudentService struct {
	uc *student.StudentUsecase
}

// NewStudentService creates a new StudentService.
func NewStudentService(uc *student.StudentUsecase) *StudentService {
	return &StudentService{uc: uc}
}
