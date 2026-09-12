package course

import "context"

// DeleteCourse 删除课程
//
// 课程不存在时报 ErrCourseNotFound。
func (h *Handler) DeleteCourse(ctx context.Context, id string) (err error) {
	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { err = h.tx.End(ctx, err) }()

	return h.CourseCmd.Delete(id)
}
