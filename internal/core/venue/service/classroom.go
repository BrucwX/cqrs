package service

import (
	"github.com/go-kratos/kratos-layout/internal/core/venue/biz/classroom"
)

// ClassroomService is a classroom service.
type ClassroomService struct {
	uc *classroom.ClassroomUsecase
}

// NewClassroomService creates a new ClassroomService.
func NewClassroomService(uc *classroom.ClassroomUsecase) *ClassroomService {
	return &ClassroomService{uc: uc}
}
