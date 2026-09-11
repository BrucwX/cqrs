package command

import (
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

// Save 保存课程（新增或更新）
func (c *CourseCommand) Save(crs *course.Course) error {
	if crs == nil {
		return repo.ErrCourseRequired
	}
	c.data.SaveCourse(crs)
	return nil
}

// Delete 删除课程
func (c *CourseCommand) Delete(id string) error {
	if !c.data.DeleteCourse(id) {
		return fmt.Errorf("%w: %s", repo.ErrCourseNotFound, id)
	}
	return nil
}
