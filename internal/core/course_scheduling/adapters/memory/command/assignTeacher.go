package command

import (
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// AssignTeacherRepo 是 repo.AssignTeacherRepo 的内存实现。
//
// 它只做「按需取数据」：排讲师的规则判定在命令服务里（checkTeacher），
// 这里不参与任何规则。
type AssignTeacherRepo struct {
	data *memory.Data
}

// 编译期断言：内存实现必须满足接口。
var _ repo.AssignTeacherRepo = (*AssignTeacherRepo)(nil)

// NewAssignTeacherRepo 创建排讲师冲突检查的数据来源。
func NewAssignTeacherRepo(d *memory.Data) repo.AssignTeacherRepo {
	return &AssignTeacherRepo{data: d}
}

// GetTeacherSlots 取该讲师现有的全部排期。
func (r *AssignTeacherRepo) GetTeacherSlots(teacherID int64) (courseSlot.CourseSlots, error) {
	out := make(courseSlot.CourseSlots, 0)
	for _, cs := range r.data.CourseSlots() {
		if cs.TeacherID() == teacherID {
			out = append(out, *cs)
		}
	}
	return out, nil
}

// GetQualifications 取该讲师持有的全部资质。
func (r *AssignTeacherRepo) GetQualifications(teacherID int64) ([]qualification.Qualification, error) {
	out := make([]qualification.Qualification, 0)
	for _, q := range r.data.Qualifications() {
		if q.TeacherID() == teacherID {
			out = append(out, *q)
		}
	}
	return out, nil
}

// GetCourseType 取某门课程归属的课程类型。
func (r *AssignTeacherRepo) GetCourseType(courseID string) (courseType.CourseType, error) {
	courses := indexBy(r.data.Courses(), func(item *course.Course) string { return item.ID() })
	item, ok := courses[courseID]
	if !ok {
		return courseType.CourseType{}, fmt.Errorf("%w: course %s", repo.ErrCourseTypeNotFound, courseID)
	}

	types := indexBy(r.data.CourseTypes(), func(item *courseType.CourseType) string { return item.ID() })
	ct, ok := types[item.CourseTypeID()]
	if !ok {
		return courseType.CourseType{}, fmt.Errorf("%w: %s", repo.ErrCourseTypeNotFound, item.CourseTypeID())
	}

	return *ct, nil
}
