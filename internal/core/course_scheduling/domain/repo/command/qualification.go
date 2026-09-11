package command

import (
	"errors"

	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
)

var (
	// ErrQualificationRequired 传入的授课资质为空。
	ErrQualificationRequired = errors.New("qualification is required")
	// ErrQualificationNotFound 指定的授课资质不存在。
	ErrQualificationNotFound = errors.New("qualification not found")
)

// QualificationCommand 授课资质命令接口
type QualificationCommand interface {
	// Save 保存授课资质（新增或更新）
	Save(q *qualification.Qualification) error
	// Delete 删除授课资质
	Delete(id int64) error
}
