package qualificationCheck

import "context"

// CheckTeacherQualified 判断该讲师有没有资格教这一门课。
//
// 先把课程解析成它所属的课程类型，再交给 CheckTeacherQualifiedForCourseType 判。
func (s *Service) CheckTeacherQualified(ctx context.Context, teacherID int64, courseID string) (bool, error) {
	cty, err := s.courseTypes.GetCourseType(courseID)
	if err != nil {
		return false, err
	}

	return s.CheckTeacherQualifiedForCourseType(ctx, teacherID, cty.ID())
}
