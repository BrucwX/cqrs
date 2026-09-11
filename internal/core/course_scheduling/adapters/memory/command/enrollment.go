package command

import (
	"context"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
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
// 先把 checkEnrollment 需要的上下文装好交给调用方判定（选课窗口 + 容量 + 时间冲突），
// 通过后才写入；传 nil 表示不做检查。
func (c *EnrollmentCommand) Enroll(ctx context.Context, e *enrollment.CourseEnrollment, checkConflictFn func(ctx context.Context, ec repo.EnrollmentContext) (bool, error)) error {
	if e == nil {
		return repo.ErrEnrollmentRequired
	}

	if checkConflictFn != nil {
		conflict, err := checkConflictFn(ctx, c.enrollmentContext(e))
		if err != nil {
			return err
		}
		if conflict {
			return fmt.Errorf("%w: student %d course %s", repo.ErrEnrollmentConflict, e.StudentID(), e.CourseID())
		}
	}

	return c.Save(e)
}

// --- 内部实现 ---

// enrollmentContext 组装选课时冲突检查所需的上下文。
func (c *EnrollmentCommand) enrollmentContext(e *enrollment.CourseEnrollment) repo.EnrollmentContext {
	return repo.EnrollmentContext{
		Course:       c.courseByID(e.CourseID()),
		TargetSlots:  c.slotsOfCourses(e.CourseID()),
		StudentSlots: c.slotsOfStudent(e.StudentID()),
	}
}

// courseByID 取课程，取不到返回零值。
func (c *EnrollmentCommand) courseByID(id string) course.Course {
	for _, item := range c.data.Courses() {
		if item.ID() == id {
			return *item
		}
	}
	return course.Course{}
}

// slotsOfCourses 取若干门课程的全部排期。
func (c *EnrollmentCommand) slotsOfCourses(courseIDs ...string) courseSlot.CourseSlots {
	want := make(map[string]struct{}, len(courseIDs))
	for _, id := range courseIDs {
		want[id] = struct{}{}
	}

	out := make(courseSlot.CourseSlots, 0)
	for _, cs := range c.data.CourseSlots() {
		if _, ok := want[cs.CourseID()]; ok {
			out = append(out, *cs)
		}
	}
	return out
}

// slotsOfStudent 取该学员「在学」课程的全部排期（含目标课程本身，
// 因此重复选课会表现为与自身槽位重叠）。
// 已退课 / 已结业的记录不计入。
func (c *EnrollmentCommand) slotsOfStudent(studentID int64) courseSlot.CourseSlots {
	courseIDs := make([]string, 0)
	for _, e := range c.data.Enrollments() {
		if e.StudentID() == studentID && e.IsActive() {
			courseIDs = append(courseIDs, e.CourseID())
		}
	}
	return c.slotsOfCourses(courseIDs...)
}
