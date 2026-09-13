package student

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
)

// StudentInput 保存学员命令
//
// ID 为 nil 表示新增（学员 ID 由聚合生成）；
// 非 nil 表示更新，指的学员不存在时报 student.ErrStudentNotFound。
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
func (h *Handler) SaveStudent(ctx context.Context, cmd StudentInput) (err error) {
	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return err
	}
	// 命名返回值 + 把 End 的结果赋回去：defer 里读到的 err 才是函数真正要返回的
	// 那个。若写成外层 := 声明的 err，defer 捕获的是 Begin 的 nil，出错也会提交。
	defer func() { err = h.tx.End(ctx, err) }()

	if cmd.ID == nil {
		created, err := student.NewStudent(cmd.Name, cmd.StudentType, cmd.Contact, cmd.Status)
		if err != nil {
			return err
		}
		return h.StudentCmd.CreateStudent(ctx, created)
	}

	return h.StudentCmd.UpdateStudent(
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
