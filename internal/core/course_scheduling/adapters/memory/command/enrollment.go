package command

import (
	"context"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// EnrollmentCommand 课程注册命令实现（内存版）。
type EnrollmentCommand struct {
	data *memory.Data
}

// 编译期断言：内存实现必须满足命令接口。
var _ repo.CourseEnrollmentCommand = (*EnrollmentCommand)(nil)

// NewEnrollmentCommand 创建课程注册命令实现
//
// 返回接口类型，这样实现不完整时会在编译期暴露。
func NewEnrollmentCommand(d *memory.Data) repo.CourseEnrollmentCommand {
	return &EnrollmentCommand{data: d}
}

// Save 保存课程注册（新增或更新）
func (c *EnrollmentCommand) Save(e *enrollment.CourseEnrollment) error {
	if e == nil {
		return repo.ErrEnrollmentRequired
	}
	c.data.SaveEnrollment(e)
	return nil
}

// Delete 删除课程注册
func (c *EnrollmentCommand) Delete(id int64) error {
	if !c.data.DeleteEnrollment(id) {
		return fmt.Errorf("%w: %d", repo.ErrEnrollmentNotFound, id)
	}
	return nil
}

// Enroll 学员选课
//
// 先把该课程聚合取出来交给调用方判定（选课窗口 + 容量 + 时间冲突），
// 通过后才写入；传 nil 表示不做检查。
func (c *EnrollmentCommand) Enroll(ctx context.Context, e *enrollment.CourseEnrollment, checkConflictFn func(ctx context.Context, e *enrollment.CourseEnrollment, crs course.Course) (bool, error)) error {
	if e == nil {
		return repo.ErrEnrollmentRequired
	}

	if checkConflictFn != nil {
		crs, err := c.courseByID(e.CourseID())
		if err != nil {
			return err
		}

		conflict, err := checkConflictFn(ctx, e, crs)
		if err != nil {
			return err
		}
		if conflict {
			return fmt.Errorf("%w: student %d course %s", repo.ErrEnrollmentConflict, e.StudentID(), e.CourseID())
		}
	}

	return c.Save(e)
}

// courseByID 取课程聚合，取不到报 not found。
func (c *EnrollmentCommand) courseByID(id string) (course.Course, error) {
	for _, item := range c.data.Courses() {
		if item.ID() == id {
			return *item, nil
		}
	}
	return course.Course{}, fmt.Errorf("%w: %s", repo.ErrCourseNotFound, id)
}
