package teacher

import (
	"time"

	"github.com/google/uuid"
)

// Teacher represents a teacher. It is the teacher aggregate root.
type Teacher struct {
	ID        string
	Name      string
	Phone     string
	Email     string
	Bio       string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// NewTeacher creates a new teacher.
func NewTeacher(name, phone, email string) *Teacher {
	return &Teacher{
		ID:        uuid.New().String(),
		Name:      name,
		Phone:     phone,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Validate validates the teacher.
func (t *Teacher) Validate() error {
	if t.Name == "" {
		return ErrTeacherInvalidArgument
	}
	if t.Phone == "" {
		return ErrTeacherInvalidArgument
	}
	return nil
}
