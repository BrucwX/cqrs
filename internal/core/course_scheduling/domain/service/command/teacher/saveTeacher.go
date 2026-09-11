package teacher

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
)

// TeacherInput 保存讲师命令
//
// ID 为 nil 表示新增（讲师 ID 由聚合生成）；
// 非 nil 表示更新，指的讲师不存在时报 repo.ErrTeacherNotFound。
// 更新时只应用非 nil 的字段，nil 的字段保持原值。
//
// StudentID 只在新增时有意义（改不了，那是身份），不给就由聚合生成。
type TeacherInput struct {
	ID        *int64
	StudentID *int64
	Name      *string
	Title     *string
	Contact   *teacher.ContactInfo
	Status    *teacher.Status
}

// SaveTeacher 保存或更新讲师
//
// 新增时 Name 必填；更新时只应用非 nil 的字段。
// 必填校验、状态机与字段变更都由聚合决定。
func (h *Handler) SaveTeacher(ctx context.Context, cmd TeacherInput) error {
	if cmd.ID == nil {
		created, err := teacher.NewTeacher(cmd.StudentID, cmd.Name, cmd.Title, cmd.Contact, cmd.Status)
		if err != nil {
			return err
		}
		return h.TeacherCmd.Create(created)
	}

	return h.TeacherCmd.Update(
		ctx,
		*cmd.ID,
		func(_ context.Context, t *teacher.Teacher) (*teacher.Teacher, error) {
			if err := t.Update(cmd.Title, cmd.Contact, cmd.Status); err != nil {
				return nil, err
			}
			return t, nil
		},
	)
}
