package courseSlot

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
)

// checkTeacherQualification 检查老师是否有能力上这门课
// 返回 true 表示老师没有资质（有冲突），false 表示老师有资质（没有冲突）
func checkTeacherQualification(ctx context.Context, cty courseType.CourseType, t_q []qualification.Qualification) (bool, error) {
	// 检查老师 ID是否匹配
	for _, q := range t_q {
		if q.CourseTypeID() == cty.ID() {
			return false, nil // 老师有资质，没有冲突
		}
	}

	return true, nil // 老师没有资质，有冲突
}
