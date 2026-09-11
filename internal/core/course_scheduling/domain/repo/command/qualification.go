package command

import (
	"context"
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
//
// 资质只「发」和「撤」：能不能拿到由领域规则决定（见 domain service 的 Qualify），
// 不存在任人填字段的通用更新。
type QualificationCommand interface {
	// GrantQualification 授予授课资质
	//
	// checkQualifiedFn 由调用方注入，仓库会把「待授予的资质」传进去；
	// 传 nil 表示不做检查。判定不通过时不写入。
	GrantQualification(
		ctx context.Context,
		q *qualification.Qualification,
		checkQualifiedFn func(ctx context.Context,
			 q *qualification.Qualification) (bool, error),
	) error
	// Delete 删除授课资质
	Delete(id int64) error
}
