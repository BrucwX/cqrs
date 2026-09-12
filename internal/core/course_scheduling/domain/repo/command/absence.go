package command

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/absence"
)

// AbsenceRecordCommand 缺勤记录命令接口
//
// 只归一类：凡是取「缺勤记录」的方法都放这里（按返回值的聚合根归类，
// 不看是哪个用例在用）。
type AbsenceRecordCommand interface {
	// Save 保存缺勤记录（新增或更新）
	Save(ctx context.Context, a *absence.AbsenceRecord) error
	// Delete 删除缺勤记录
	Delete(ctx context.Context, id int64) error
	// MustGet 取缺勤记录聚合；不存在时报 ErrAbsenceNotFound
	MustGet(ctx context.Context, id int64) (absence.AbsenceRecord, error)
	// GetAbsences 取该学员的全部缺勤记录
	GetAbsences(ctx context.Context, studentID int64) ([]absence.AbsenceRecord, error)
}
