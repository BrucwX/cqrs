package command

import (
	"errors"

	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
)

var (
	// ErrAbsenceRequired 传入的缺勤记录为空。
	ErrAbsenceRequired = errors.New("absence record is required")
	// ErrAbsenceNotFound 指定的缺勤记录不存在。
	ErrAbsenceNotFound = errors.New("absence record not found")
)

// AbsenceRecordCommand 缺勤记录命令接口
//
// 只归一类：凡是取「缺勤记录」的方法都放这里（按返回值的聚合根归类，
// 不看是哪个用例在用）。
type AbsenceRecordCommand interface {
	// Save 保存缺勤记录（新增或更新）
	Save(a *absence.AbsenceRecord) error
	// Delete 删除缺勤记录
	Delete(id int64) error
	// GetAbsences 取该学员的全部缺勤记录
	GetAbsences(studentID int64) ([]absence.AbsenceRecord, error)
}
