package course

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
)

// CourseInput 保存课程命令
//
// ID 为 nil 表示新增（课程 ID 由聚合生成）；
// 非 nil 表示更新，指的课程不存在时报 course.ErrCourseNotFound。
// 更新时只应用非 nil 的字段，nil 的字段保持原值。
type CourseInput struct {
	ID           *string
	CourseTypeID *string
	Capacity     *course.Capacity
	Enrollment   *course.EnrollmentWindow
	Period       *course.CoursePeriod
}

// SaveCourse 保存或更新课程
//
// 新增时 CourseTypeID / Capacity / Enrollment / Period 必填；更新时只应用非 nil 的字段。
// 必填校验与字段变更都由聚合决定。
func (h *Handler) SaveCourse(ctx context.Context, cmd CourseInput) (err error) {
	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { err = h.tx.End(ctx, err) }()

	if cmd.ID == nil {
		created, err := course.NewCourse(cmd.CourseTypeID, cmd.Capacity, cmd.Enrollment, cmd.Period)
		if err != nil {
			return err
		}
		return h.CourseCmd.CreateCourse(ctx, created)
	}

	return h.CourseCmd.UpdateCourse(
		ctx,
		*cmd.ID,
		func(_ context.Context, crs *course.Course) (*course.Course, error) {
			if err := crs.Update(cmd.CourseTypeID, cmd.Capacity, cmd.Enrollment, cmd.Period); err != nil {
				return nil, err
			}
			return crs, nil
		},
	)
}
