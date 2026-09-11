package qualification

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
)

// QualifyInput 授予授课资质命令
type QualifyInput struct {
	TeacherID    int64
	CourseTypeID string
}

// Qualify 授予授课资质
//
// 仓库（命令适配器）会把「待写入的资质」传进检查里，判定走 h.checkFinished；
// 判定与写证是同一次调用，不够格就不会落库。
func (h *Handler) Qualify(ctx context.Context, cmd QualifyInput) error {
	created, err := qualification.NewQualification(cmd.TeacherID, cmd.CourseTypeID)
	if err != nil {
		return err
	}

	return h.QualificationCmd.GrantQualification(ctx, created, h.checkFinished)
}

// checkFinished 判断这个讲师够不够格发证。
//
// 返回 true 表示没修完（拒绝发证），false 表示可以发。
// 规则 = 该类型下有某门课，讲师在那门课上已结业，且那门课没有任何缺勤记录。
// 缺勤按课程算：在同一类型的另一门课上缺过课，不影响这一门的判定。
//
// 讲师是以「学员」身份参训的，所以报名/缺勤记录挂的是 t.StudentID() 而不是讲师 ID。
// 仓库只把「待写入的资质」传进来，课程/报名/缺勤按需自己取。
func (h *Handler) checkFinished(ctx context.Context, q *qualification.Qualification) (bool, error) {
	t, err := h.qualify.GetTeacher(q.TeacherID())
	if err != nil {
		return false, err
	}

	courses, err := h.qualify.GetCourses(q.CourseTypeID())
	if err != nil {
		return false, err
	}

	enrollments, err := h.qualify.GetEnrollments(t.StudentID())
	if err != nil {
		return false, err
	}

	absences, err := h.qualify.GetAbsences(t.StudentID())
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
		return false, nil
	}

	return true, nil
}
