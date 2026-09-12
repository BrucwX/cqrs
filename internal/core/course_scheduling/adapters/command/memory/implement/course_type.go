package implement

import (
	"context"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/command/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// CourseTypeCommand 课程类型命令实现（内存版）。
//
// 课程类型目前没有增删改的用例，这里只有「按课程取它归属的类型」——
// 排讲师核对资质、授资质判定都要用它。
type CourseTypeCommand struct {
	data *memory.Data
}

// 编译期断言：内存实现必须满足命令接口。
var _ repo.CourseTypeCommand = (*CourseTypeCommand)(nil)

// NewCourseTypeCommand 创建课程类型命令实现。
func NewCourseTypeCommand(d *memory.Data) repo.CourseTypeCommand {
	return &CourseTypeCommand{data: d}
}

// GetCourseType 取某门课程归属的课程类型。
func (c *CourseTypeCommand) GetCourseType(ctx context.Context, courseID string) (courseType.CourseType, error) {
	courses := indexBy(c.data.Courses(), func(item *course.Course) string { return item.ID() })
	item, ok := courses[courseID]
	if !ok {
		return courseType.CourseType{}, fmt.Errorf("%w: course %s", repo.ErrCourseTypeNotFound, courseID)
	}

	types := indexBy(c.data.CourseTypes(), func(item *courseType.CourseType) string { return item.ID() })
	ct, ok := types[item.CourseTypeID()]
	if !ok {
		return courseType.CourseType{}, fmt.Errorf("%w: %s", repo.ErrCourseTypeNotFound, item.CourseTypeID())
	}

	return *ct, nil
}
