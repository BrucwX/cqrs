package command

import (
	"context"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// QualificationCommand 授课资质命令实现（内存版）。
type QualificationCommand struct {
	data *memory.Data
}

// 编译期断言：内存实现必须满足命令接口。
var _ repo.QualificationCommand = (*QualificationCommand)(nil)

// NewQualificationCommand 创建授课资质命令实现
//
// 返回接口类型，这样实现不完整时会在编译期暴露。
func NewQualificationCommand(d *memory.Data) repo.QualificationCommand {
	return &QualificationCommand{data: d}
}

// GrantQualification 授予授课资质
//
// 先把待授予的资质交给调用方判定（讲师修完课程了没），通过后才写入；
// 传 nil 表示不做检查。
func (c *QualificationCommand) GrantQualification(ctx context.Context, q *qualification.Qualification, checkQualifiedFn func(ctx context.Context, q *qualification.Qualification) (bool, error)) error {
	if q == nil {
		return repo.ErrQualificationRequired
	}

	if checkQualifiedFn != nil {
		notQualified, err := checkQualifiedFn(ctx, q)
		if err != nil {
			return err
		}
		if notQualified {
			return fmt.Errorf("%w: teacher %d courseType %s", repo.ErrCourseNotFinished, q.TeacherID(), q.CourseTypeID())
		}
	}

	c.data.SaveQualification(q)
	return nil
}

// Delete 删除授课资质
func (c *QualificationCommand) Delete(id int64) error {
	if !c.data.DeleteQualification(id) {
		return fmt.Errorf("%w: %d", repo.ErrQualificationNotFound, id)
	}
	return nil
}
