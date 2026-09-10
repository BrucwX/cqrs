package command

import "cqrs/internal/core/course_scheduling/domain/aggregate/absence"

// AbsenceRecordCommand 缺勤记录命令接口
type AbsenceRecordCommand interface {
	// Save 保存缺勤记录（新增或更新）
	Save(a *absence.AbsenceRecord) error
	// Delete 删除缺勤记录
	Delete(id int64) error
}
