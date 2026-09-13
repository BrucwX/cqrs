package implement

import (
	"context"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/memImp4test/command"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// CourseCommand 课程命令实现（内存版）。
type CourseCommand struct {
	data *command.Data
}

// 编译期断言：内存实现必须满足命令接口。
var _ repo.CourseCommand = (*CourseCommand)(nil)

// NewCourseCommand 创建课程命令实现
//
// 返回接口类型，这样实现不完整时会在编译期暴露。
func NewCourseCommand(d *command.Data) repo.CourseCommand {
	return &CourseCommand{data: d}
}

// CreateCourse 新增课程
func (c *CourseCommand) CreateCourse(ctx context.Context, crs *course.Course) error {
	if crs == nil {
		return course.ErrCourseRequired
	}
	c.data.SaveCourse(crs)
	return nil
}

// UpdateCourse 按 ID 取出课程交给 updateFn 改，改完写回
func (c *CourseCommand) UpdateCourse(ctx context.Context, id string, updateFn func(ctx context.Context, crs *course.Course) (*course.Course, error)) error {
	current, err := c.GetCourse(ctx, id)
	if err != nil {
		return err
	}
	if current == nil {
		return fmt.Errorf("%w: %s", course.ErrCourseNotFound, id)
	}

	updated, err := updateFn(ctx, current)
	if err != nil {
		return err
	}
	if updated == nil {
		return course.ErrCourseRequired
	}

	c.data.SaveCourse(updated)
	return nil
}

// DeleteCourse 删除课程
func (c *CourseCommand) DeleteCourse(ctx context.Context, id string) error {
	if !c.data.DeleteCourse(id) {
		return fmt.Errorf("%w: %s", course.ErrCourseNotFound, id)
	}
	return nil
}

// GetCourse 取课程；不存在时返回 (nil, nil)。
func (c *CourseCommand) GetCourse(ctx context.Context, id string) (*course.Course, error) {
	item, ok := c.data.CourseByID(id)
	if !ok {
		return nil, nil
	}
	return item, nil
}

// MustGetCourse 取课程聚合本身；取不到报 ErrCourseNotFound。
//
// 与 Get 的差别：Get 找不到时是 (nil, nil)（给「查到了没」的调用方），
// MustGetCourse 是规则判定要用的，找不到必须报错。
func (c *CourseCommand) MustGetCourse(ctx context.Context, id string) (course.Course, error) {
	item, ok := c.data.CourseByID(id)
	if !ok {
		return course.Course{}, fmt.Errorf("%w: %s", course.ErrCourseNotFound, id)
	}
	return *item, nil
}

// GetCourses 取某课程类型下的全部课程。
func (c *CourseCommand) GetCourses(ctx context.Context, courseTypeID string) ([]course.Course, error) {
	out := make([]course.Course, 0)
	for _, item := range c.data.Courses() {
		if item.CourseTypeID() == courseTypeID {
			out = append(out, *item)
		}
	}
	return out, nil
}
