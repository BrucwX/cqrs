package implement

import (
	"context"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/memImp4test/command"
	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// AbsenceCommand 缺勤记录命令实现（内存版）。
type AbsenceCommand struct {
	data *command.Data
}

// 编译期断言：内存实现必须满足命令接口。
var _ repo.AbsenceRecordCommand = (*AbsenceCommand)(nil)

// NewAbsenceCommand 创建缺勤记录命令实现
//
// 返回接口类型，这样实现不完整时会在编译期暴露。
func NewAbsenceCommand(d *command.Data) repo.AbsenceRecordCommand {
	return &AbsenceCommand{data: d}
}

// SaveAbsence 保存缺勤记录（新增或更新）
func (c *AbsenceCommand) SaveAbsence(ctx context.Context, a *absence.AbsenceRecord) error {
	if a == nil {
		return absence.ErrAbsenceRequired
	}
	c.data.SaveAbsence(a)
	return nil
}

// DeleteAbsence 删除缺勤记录
func (c *AbsenceCommand) DeleteAbsence(ctx context.Context, id int64) error {
	if !c.data.DeleteAbsence(id) {
		return fmt.Errorf("%w: %d", absence.ErrAbsenceNotFound, id)
	}
	return nil
}

// MustGetAbsence 取缺勤记录聚合；不存在时报 ErrAbsenceNotFound。
func (c *AbsenceCommand) MustGetAbsence(ctx context.Context, id int64) (absence.AbsenceRecord, error) {
	for _, item := range c.data.Absences() {
		if item.ID() == id {
			return *item, nil
		}
	}
	return absence.AbsenceRecord{}, fmt.Errorf("%w: %d", absence.ErrAbsenceNotFound, id)
}

// GetAbsences 取该学员的全部缺勤记录。
func (c *AbsenceCommand) GetAbsences(ctx context.Context, studentID int64) ([]absence.AbsenceRecord, error) {
	out := make([]absence.AbsenceRecord, 0)
	for _, item := range c.data.Absences() {
		if item.StudentID() == studentID {
			out = append(out, *item)
		}
	}
	return out, nil
}
