package student

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
)

// StudentInput 保存学员命令
//
// ID 为 nil 表示新增（学员 ID 由聚合生成）；
// 非 nil 表示更新，指的学员不存在时报 repo.ErrStudentNotFound。
// 更新时只应用非 nil 的字段，nil 的字段保持原值。
type StudentInput struct {
	ID          *int64
	Name        *string
	StudentType *student.StudentType
	Contact     *student.ContactInfo
	Status      *student.Status
}

// SaveStudent 保存或更新学员
//
// 新增时 Name 必填；更新时只应用非 nil 的字段。
// 必填校验、状态机与字段变更都由聚合决定。
func (h *Handler) SaveStudent(ctx context.Context, cmd StudentInput) error {
	if cmd.ID == nil {
		created, err := student.NewStudent(cmd.Name, cmd.StudentType, cmd.Contact, cmd.Status)
		if err != nil {
			return err
		}
		return h.StudentCmd.Create(created)
	}

	return h.StudentCmd.Update(
		ctx,
		*cmd.ID,
		func(_ context.Context, s *student.Student) (*student.Student, error) {
			if err := s.Update(cmd.Contact, cmd.Status); err != nil {
				return nil, err
			}
			return s, nil
		},
	)
}
