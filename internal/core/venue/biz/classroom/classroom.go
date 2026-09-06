package classroom

import (
	"time"

	"github.com/google/uuid"
)

// Classroom represents a classroom. It is the classroom aggregate root.
type Classroom struct {
	ID        string
	Name      string
	Building  string
	Floor     int32
	Capacity  int32
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// NewClassroom creates a new classroom.
func NewClassroom(name string, capacity int32) *Classroom {
	return &Classroom{
		ID:       uuid.New().String(),
		Name:     name,
		Capacity: capacity,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Validate validates the classroom.
func (c *Classroom) Validate() error {
	if c.Name == "" {
		return ErrClassroomInvalidArgument
	}
	if c.Capacity <= 0 {
		return ErrClassroomInvalidArgument
	}
	return nil
}
