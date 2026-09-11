package course

import "context"

// DeleteCourse 删除课程
//
// 课程不存在时报 ErrCourseNotFound。
func (h *Handler) DeleteCourse(ctx context.Context, id string) error {
	return h.CourseCmd.Delete(id)
}
