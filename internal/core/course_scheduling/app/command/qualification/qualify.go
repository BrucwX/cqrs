package qualification

import (
	"context"
	"fmt"

	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
)

// QualifyInput 授予授课资质命令
type QualifyInput struct {
	TeacherID    int64
	CourseTypeID string
}

// Qualify 授予授课资质
//
// 顺序：先判、过了才写。判定交给领域服务 QualificationCheck.CheckTeacherGrantable，
// 写回走 QualificationCmd.GrantQualification。
func (h *Handler) Qualify(ctx context.Context, cmd QualifyInput) (err error) {
	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { err = h.tx.End(ctx, err) }()

	created, err := qualification.NewQualification(cmd.TeacherID, cmd.CourseTypeID)
	if err != nil {
		return err
	}

	notFinished, err := h.qualify.CheckTeacherGrantable(ctx, cmd.TeacherID, cmd.CourseTypeID)
	if err != nil {
		return err
	}
	if notFinished {
		return fmt.Errorf("%w: teacher %d courseType %s", qualification.ErrCourseNotFinished, cmd.TeacherID, cmd.CourseTypeID)
	}

	return h.QualificationCmd.GrantQualification(ctx, created)
}
