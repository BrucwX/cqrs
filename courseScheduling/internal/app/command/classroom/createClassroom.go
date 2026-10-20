package classroom

import (
	"context"

	"courseScheduling/internal/domain/aggregate/classroom"
)

type CreateClassroomInput struct {
	Location *classroom.Location
	Capacity *int
	Status   *classroom.Status
}

func (h *Handler) CreateClassroom(ctx context.Context, cmd CreateClassroomInput) (err error) {
	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { err = h.tx.End(ctx, err) }()

	created, err := classroom.NewClassroom(cmd.Location, cmd.Capacity, cmd.Status)
	if err != nil {
		return err
	}

	return h.ClassroomCmd.CreateClassroom(ctx, created)
}
