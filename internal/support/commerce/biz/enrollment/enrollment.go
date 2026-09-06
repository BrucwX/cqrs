package enrollment

import (
	"time"

	"github.com/google/uuid"
)

// EnrollmentStatus represents the enrollment status.
type EnrollmentStatus string

const (
	EnrollmentStatusPending  EnrollmentStatus = "PENDING"
	EnrollmentStatusActive   EnrollmentStatus = "ACTIVE"
	EnrollmentStatusCanceled EnrollmentStatus = "CANCELED"
)

// Enrollment represents a student enrollment in a course. It is the enrollment aggregate root.
type Enrollment struct {
	ID        string
	StudentID string
	CourseID  string
	Status    EnrollmentStatus
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// NewEnrollment creates a new enrollment.
func NewEnrollment(studentID, courseID string) *Enrollment {
	return &Enrollment{
		ID:        uuid.New().String(),
		StudentID: studentID,
		CourseID:  courseID,
		Status:    EnrollmentStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Activate activates the enrollment.
func (e *Enrollment) Activate() {
	e.Status = EnrollmentStatusActive
	e.UpdatedAt = time.Now()
}

// Cancel cancels the enrollment.
func (e *Enrollment) Cancel() {
	e.Status = EnrollmentStatusCanceled
	e.UpdatedAt = time.Now()
}

// IsActive checks if the enrollment is active.
func (e *Enrollment) IsActive() bool {
	return e.Status == EnrollmentStatusActive
}
