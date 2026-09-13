package implement

import (
	"context"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/test_memory/command"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// QualificationCommand 授课资质命令实现（内存版）。
type QualificationCommand struct {
	data *command.Data
}

// 编译期断言：内存实现必须满足命令接口。
var _ repo.QualificationCommand = (*QualificationCommand)(nil)

// NewQualificationCommand 创建授课资质命令实现
//
// 返回接口类型，这样实现不完整时会在编译期暴露。
func NewQualificationCommand(d *command.Data) repo.QualificationCommand {
	return &QualificationCommand{data: d}
}

// GrantQualification 授予授课资质
//
// 只管写：够不够格（讲师修没修完该类型的课）由调用方在调过来之前判完。
func (c *QualificationCommand) GrantQualification(ctx context.Context, q *qualification.Qualification) error {
	if q == nil {
		return qualification.ErrQualificationRequired
	}

	c.data.SaveQualification(q)
	return nil
}

// DeleteQualification 删除授课资质
func (c *QualificationCommand) DeleteQualification(ctx context.Context, id int64) error {
	if !c.data.DeleteQualification(id) {
		return fmt.Errorf("%w: %d", qualification.ErrQualificationNotFound, id)
	}
	return nil
}

// MustGetQualification 取授课资质聚合；不存在时报 ErrQualificationNotFound。
func (c *QualificationCommand) MustGetQualification(ctx context.Context, id int64) (qualification.Qualification, error) {
	for _, item := range c.data.Qualifications() {
		if item.ID() == id {
			return *item, nil
		}
	}
	return qualification.Qualification{}, fmt.Errorf("%w: %d", qualification.ErrQualificationNotFound, id)
}

// GetQualifications 取该讲师持有的全部资质。
func (c *QualificationCommand) GetQualifications(ctx context.Context, teacherID int64) ([]qualification.Qualification, error) {
	out := make([]qualification.Qualification, 0)
	for _, item := range c.data.Qualifications() {
		if item.TeacherID() == teacherID {
			out = append(out, *item)
		}
	}
	return out, nil
}
