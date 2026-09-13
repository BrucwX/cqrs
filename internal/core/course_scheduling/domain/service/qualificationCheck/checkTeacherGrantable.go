package qualificationCheck

import "context"

// CheckTeacherGrantable 判断该不该给这位讲师发这个课程类型的资质。
//
// 规则 = 该课程类型下有某门课，讲师在那门课上已结业，且那门课没有任何缺勤记录。
// 缺勤按课程算：在同一类型的另一门课上缺过课，不影响这一门的判定。
//
// 与 CheckTeacherQualifiedForCourseType 正好相反：那边问「他手上已经有什么资质」
// （能不能排课），这边问「该不该给他发这张证」。
//
// 讲师是以「学员」身份参训的，所以报名 / 缺勤记录挂的是 t.StudentID() 而不是
// 讲师 ID —— 这层换算在这里做，调用方手上只有讲师 ID。
//
// 讲师不存在时报 not found —— 调用方传错了 ID 应该报出来，
// 而不是当成「不够格」悄悄拦下。
func (s *Service) CheckTeacherGrantable(ctx context.Context, teacherID int64, courseTypeID string) (bool, error) {
	t, err := s.teachers.MustGetTeacher(ctx, teacherID)
	if err != nil {
		return false, err
	}

	courses, err := s.courses.GetCourses(ctx, courseTypeID)
	if err != nil {
		return false, err
	}

	enrollments, err := s.enrollments.GetEnrollments(ctx, t.StudentID())
	if err != nil {
		return false, err
	}

	absences, err := s.absences.GetAbsences(ctx, t.StudentID())
	if err != nil {
		return false, err
	}

	completed := make(map[string]struct{}, len(enrollments))
	for _, e := range enrollments {
		if e.IsCompleted() {
			completed[e.CourseID()] = struct{}{}
		}
	}

	absent := make(map[string]struct{}, len(absences))
	for _, a := range absences {
		absent[a.CourseID()] = struct{}{}
	}

	for _, c := range courses {
		if _, ok := completed[c.ID()]; !ok {
			continue // 这门课没结业，换下一门看
		}
		if _, ok := absent[c.ID()]; ok {
			continue // 这门课缺过课，不算修完
		}
		return false, nil // 找到一门修完的了 -> 够格
	}

	return true, nil
}
