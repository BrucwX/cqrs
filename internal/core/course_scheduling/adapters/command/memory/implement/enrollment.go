package implement

import (
	"context"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/command/memory"
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
func (c *EnrollmentCommand) Save(ctx context.Context, e *enrollment.CourseEnrollment) error {
	if e == nil {
		return repo.ErrEnrollmentRequired
	}
	c.data.SaveEnrollment(e)
	return nil
}

// Delete 删除课程注册
func (c *EnrollmentCommand) Delete(ctx context.Context, id int64) error {
	if !c.data.DeleteEnrollment(id) {
		return fmt.Errorf("%w: %d", repo.ErrEnrollmentNotFound, id)
	}
	return nil
}

// Enroll 学员选课
//
// 只管写：准入判定（选课窗口 / 容量 / 时间冲突）由调用方在调过来之前做完。
func (c *EnrollmentCommand) Enroll(ctx context.Context, e *enrollment.CourseEnrollment) error {
	if e == nil {
		return repo.ErrEnrollmentRequired
	}

	return c.Save(ctx, e)
}

// GetEnrollments 取该学员的全部报名记录。
func (c *EnrollmentCommand) GetEnrollments(ctx context.Context, studentID int64) ([]enrollment.CourseEnrollment, error) {
	out := make([]enrollment.CourseEnrollment, 0)
	for _, item := range c.data.Enrollments() {
		if item.StudentID() == studentID {
			out = append(out, *item)
		}
	}
	return out, nil
}
