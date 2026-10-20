package classroom

import (
	"context"

	"courseScheduling/internal/domain/aggregate/classroom"
)

type UpdateClassroomInput struct {
	ID       string
	Location *classroom.Location
	Capacity *int
	Status   *classroom.Status
}

// UpdateClassroom 更新教室
//
// 走乐观锁事务（两阶段，见 appRepo.OptimisticTransaction）：这里只表达「怎么改」，
// 读教室、校验 lock_version、冲突重试都在适配器里；updateFn 可能被重跑。
func (h *Handler) UpdateClassroom(ctx context.Context, cmd UpdateClassroomInput) error {

	fn := func(_ context.Context, cl *classroom.Classroom) (*classroom.Classroom, error) {
		if err := cl.Update(cmd.Location, cmd.Capacity, cmd.Status); err != nil {
			return nil, err
		}
		return cl, nil
	}

	return h.optTx.RunUpdatClassroom(ctx, cmd.ID, fn)
}
