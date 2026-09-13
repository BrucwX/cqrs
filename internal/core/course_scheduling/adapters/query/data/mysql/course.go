package mysql

import (
	"context"

	"cqrs/internal/core/course_scheduling/adapters/query/data/mysql/help"
	"cqrs/internal/core/course_scheduling/adapters/query/data/mysql/model"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/enrollment"
)

// Page 分页查询课程列表，按 ID 升序。
func (q *Data) PageCourses(ctx context.Context, page, pageSize int) ([]*course.Course, error) {
	limit, offset := help.LimitOffset(page, pageSize)
	return help.QueryAll(ctx, q.Conn(ctx), `
SELECT `+help.CourseColumns+`
  FROM course
 ORDER BY id
 LIMIT ? OFFSET ?`, []any{limit, offset}, help.ScanCourse, model.CourseToDO)
}

// AvailableForStudent 返回与学员「当前在学课程」不冲突的课程。
//
// 占用集 = 该学员所有生效中（status = 在读）报名记录对应课程的全部周排期槽位。
func (q *Data) AvailableForStudent(ctx context.Context, studentID int64) ([]*course.Course, error) {
	const busy = `b.course_id IN (
                SELECT e.course_id
                  FROM course_enrollment e
                 WHERE e.student_id = ? AND e.status = ?)`

	return q.available(ctx, busy, studentID, int(enrollment.StatusEnrolled))
}

// AvailableForTeacher 返回与讲师「现有排课」不冲突的课程。
//
// 占用集 = 该讲师已被指派（teacher_id 命中）的全部周排期槽位。
func (q *Data) AvailableForTeacher(ctx context.Context, teacherID int64) ([]*course.Course, error) {
	return q.available(ctx, "b.teacher_id = ?", teacherID)
}

// AvailableForClassroom 返回与教室「现有排课」不冲突的课程。
//
// 占用集 = 该教室已被占用（classroom_id 命中）的全部周排期槽位。
func (q *Data) AvailableForClassroom(ctx context.Context, classroomID string) ([]*course.Course, error) {
	return q.available(ctx, "b.classroom_id = ?", classroomID)
}

// available 返回所有排期都不与「已占用时间」冲突的课程。
//
// busy 是描述占用集的 SQL 片段（学员看课程、讲师/教室看自己名下的槽位），
// 冲突判定 = 同一星期几 且 日内时间区间重叠：
//
//	s.start_time < b.end_time AND b.start_time < s.end_time
//
// 这里刻意不看讲师/教室：对学员而言，两门课时间撞了就是上不了。
// 没有排期槽位的课程不会被判定为冲突，因此会出现在结果里。
func (q *Data) available(ctx context.Context, busy string, args ...any) ([]*course.Course, error) {
	query := `
SELECT ` + help.Qualify(help.CourseColumns, "c") + `
  FROM course c
 WHERE NOT EXISTS (
       SELECT 1
         FROM course_slot s
         JOIN course_slot b ON b.weekday = s.weekday
                           AND s.start_time < b.end_time
                           AND b.start_time < s.end_time
        WHERE s.course_id = c.id
          AND (` + busy + `)
 )
 ORDER BY c.id`

	return help.QueryAll(ctx, q.Conn(ctx), query, args, help.ScanCourse, model.CourseToDO)
}
