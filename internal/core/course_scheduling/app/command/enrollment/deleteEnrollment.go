package enrollment

import "context"

// DeleteEnrollment 删除课程注册记录
//
// 记录不存在时报 ErrEnrollmentNotFound。注意这是硬删除，
// 需要保留退课痕迹请改用 CourseEnrollment.Drop。
func (h *Handler) DeleteEnrollment(ctx context.Context, id int64) (err error) {
	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { err = h.tx.End(ctx, err) }()

	return h.EnrollmentCmd.Delete(id)
}
