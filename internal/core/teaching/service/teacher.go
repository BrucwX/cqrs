package service

import (
	"github.com/go-kratos/kratos-layout/internal/core/teaching/biz/teacher"
)

// TeacherService is a teacher service.
type TeacherService struct {
	uc *teacher.TeacherUsecase
}

// NewTeacherService creates a new TeacherService.
func NewTeacherService(uc *teacher.TeacherUsecase) *TeacherService {
	return &TeacherService{uc: uc}
}
