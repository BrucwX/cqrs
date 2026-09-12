package classroom

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
)

// ClassroomInput 保存教室命令
//
// ID 为 nil 表示新增（教室 ID 由服务端生成）；
// 非 nil 表示更新，指的教室不存在时报 repo.ErrClassroomNotFound。
// 更新时只应用非 nil 的字段，nil 的字段保持原值。
type ClassroomInput struct {
	ID       *string
	Location *classroom.Location
	Capacity *int
	Status   *classroom.Status
}

// SaveClassroom 保存或更新教室
//
// ID 为 nil 表示新增（ID 由聚合生成）；非 nil 表示更新，教室不存在时报
// repo.ErrClassroomNotFound。必填校验、状态机、哪些字段可改都由聚合决定。
func (h *Handler) SaveClassroom(ctx context.Context, cmd ClassroomInput) (err error) {
	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { err = h.tx.End(ctx, err) }()

	if cmd.ID == nil {
		created, err := classroom.NewClassroom(cmd.Location, cmd.Capacity, cmd.Status)
		if err != nil {
			return err
		}
		return h.ClassroomCmd.Create(created)
	}

	return h.ClassroomCmd.Update(
		ctx,
		*cmd.ID,
		func(_ context.Context, cl *classroom.Classroom) (*classroom.Classroom, error) {
			if err := cl.Update(cmd.Location, cmd.Capacity, cmd.Status); err != nil {
				return nil, err
			}
			return cl, nil
		},
	)
}
