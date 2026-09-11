package enrollment

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
)

// StudentEnroll 学生选课命令
type StudentEnroll struct {
	StudentID int64
	CourseID  string
}

// StudentEnroll 学生选课
func (h *Handler) StudentEnroll(ctx context.Context, cmd StudentEnroll) (*enrollment.CourseEnrollment, error) {
	// 创建选课记录
	enroll, err := enrollment.NewCourseEnrollment(cmd.StudentID, cmd.CourseID)
	if err != nil {
		return nil, err
	}

	// 保存
	if err := h.EnrollmentCmd.Save(enroll); err != nil {
		return nil, err
	}

	return enroll, nil
}
