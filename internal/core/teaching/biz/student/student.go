package student

import (
	"time"

	"github.com/google/uuid"
)

// Student represents a student. It is the student aggregate root.
type Student struct {
	ID        string
	Name      string
	Phone     string
	Email     string
	ParentID  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// NewStudent creates a new student.
func NewStudent(name, phone, email string) *Student {
	return &Student{
		ID:        uuid.New().String(),
		Name:      name,
		Phone:     phone,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Validate validates the student.
func (s *Student) Validate() error {
	if s.Name == "" {
		return ErrStudentInvalidArgument
	}
	if s.Phone == "" {
		return ErrStudentInvalidArgument
	}
	return nil
}
