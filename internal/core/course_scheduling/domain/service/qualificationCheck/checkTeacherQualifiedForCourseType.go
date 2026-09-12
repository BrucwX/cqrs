package qualificationCheck

import "context"

// CheckTeacherQualifiedForCourseType 判断该讲师有没有资格教这个课程类型。
//
// 这是「有没有资质」这条唯一真正要算的判定 —— 资质本身就绑在「讲师 + 课程类型」
// 上，另外两个入口只是先把它要的课程类型解析出来。
//
// 讲师不存在时报 not found —— 调用方传错了 ID 应该报出来，
// 而不是当成「没资质」悄悄拦下。
func (s *Service) CheckTeacherQualifiedForCourseType(ctx context.Context, teacherID int64, courseTypeID string) (bool, error) {
	t, err := s.teachers.GetTeacher(teacherID)
	if err != nil {
		return false, err
	}

	held, err := s.qualifications.GetQualifications(t.ID())
	if err != nil {
		return false, err
	}

	for _, q := range held {
		if q.CourseTypeID() == courseTypeID {
			return false, nil // 手上有这条 -> 有资质
		}
	}

	return true, nil // 一条都不覆盖 -> 没资质
}
