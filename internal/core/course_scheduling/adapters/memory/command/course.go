package command

import (
	"context"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// CourseCommand 课程命令实现（内存版）。
type CourseCommand struct {
	data *memory.Data
}

// 编译期断言：内存实现必须满足命令接口。
var _ repo.CourseCommand = (*CourseCommand)(nil)

// NewCourseCommand 创建课程命令实现
//
// 返回接口类型，这样实现不完整时会在编译期暴露。
func NewCourseCommand(d *memory.Data) repo.CourseCommand {
	return &CourseCommand{data: d}
}

// Create 新增课程
func (c *CourseCommand) Create(crs *course.Course) error {
	if crs == nil {
		return repo.ErrCourseRequired
	}
	c.data.SaveCourse(crs)
	return nil
}

// Update 按 ID 取出课程交给 updateFn 改，改完写回
func (c *CourseCommand) Update(ctx context.Context, id string, updateFn func(ctx context.Context, crs *course.Course) (*course.Course, error)) error {
	current, err := c.Get(id)
	if err != nil {
		return err
	}
	if current == nil {
		return fmt.Errorf("%w: %s", repo.ErrCourseNotFound, id)
	}

	updated, err := updateFn(ctx, current)
	if err != nil {
		return err
	}
	if updated == nil {
		return repo.ErrCourseRequired
	}

	c.data.SaveCourse(updated)
	return nil
}

// Delete 删除课程
func (c *CourseCommand) Delete(id string) error {
	if !c.data.DeleteCourse(id) {
		return fmt.Errorf("%w: %s", repo.ErrCourseNotFound, id)
	}
	return nil
}

// Get 取课程；不存在时返回 (nil, nil)。
func (c *CourseCommand) Get(id string) (*course.Course, error) {
	item, ok := c.data.CourseByID(id)
	if !ok {
		return nil, nil
	}
	return item, nil
}
