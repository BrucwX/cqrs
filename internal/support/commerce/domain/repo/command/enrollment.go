package command

import "cqrs/internal/support/commerce/domain/aggregate/enrollment"

// EnrollmentCommand 注册命令接口
type EnrollmentCommand interface {
	// Save 保存注册（新增或更新）
	Save(e *enrollment.Enrollment) error
	// Delete 删除注册
	Delete(id string) error
}
