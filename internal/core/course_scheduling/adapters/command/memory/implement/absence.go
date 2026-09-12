package implement

import (
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/command/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// AbsenceCommand 缺勤记录命令实现（内存版）。
type AbsenceCommand struct {
	data *memory.Data
}

// 编译期断言：内存实现必须满足命令接口。
var _ repo.AbsenceRecordCommand = (*AbsenceCommand)(nil)

// NewAbsenceCommand 创建缺勤记录命令实现
//
// 返回接口类型，这样实现不完整时会在编译期暴露。
func NewAbsenceCommand(d *memory.Data) repo.AbsenceRecordCommand {
	return &AbsenceCommand{data: d}
}

// Save 保存缺勤记录（新增或更新）
func (c *AbsenceCommand) Save(a *absence.AbsenceRecord) error {
	if a == nil {
		return repo.ErrAbsenceRequired
	}
	c.data.SaveAbsence(a)
	return nil
}

// Delete 删除缺勤记录
func (c *AbsenceCommand) Delete(id int64) error {
	if !c.data.DeleteAbsence(id) {
		return fmt.Errorf("%w: %d", repo.ErrAbsenceNotFound, id)
	}
	return nil
}

// GetAbsences 取该学员的全部缺勤记录。
func (c *AbsenceCommand) GetAbsences(studentID int64) ([]absence.AbsenceRecord, error) {
	out := make([]absence.AbsenceRecord, 0)
	for _, item := range c.data.Absences() {
		if item.StudentID() == studentID {
			out = append(out, *item)
		}
	}
	return out, nil
}
