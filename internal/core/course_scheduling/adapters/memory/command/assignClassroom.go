package command

import (
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// AssignClassroomRepo 是 repo.AssignClassroomRepo 的内存实现。
//
// 它只做「按需取数据」：排教室的规则判定在命令服务里（checkClassroom），
// 这里不参与任何规则。
type AssignClassroomRepo struct {
	data *memory.Data
}

// 编译期断言：内存实现必须满足接口。
var _ repo.AssignClassroomRepo = (*AssignClassroomRepo)(nil)

// NewAssignClassroomRepo 创建排教室冲突检查的数据来源。
func NewAssignClassroomRepo(d *memory.Data) repo.AssignClassroomRepo {
	return &AssignClassroomRepo{data: d}
}

// GetClassroomSlots 取该教室现有的全部排期。
func (r *AssignClassroomRepo) GetClassroomSlots(classroomID string) (courseSlot.CourseSlots, error) {
	out := make(courseSlot.CourseSlots, 0)
	for _, cs := range r.data.CourseSlots() {
		if cs.ClassroomID() == classroomID {
			out = append(out, *cs)
		}
	}
	return out, nil
}

// GetCourse 取课程本身。
func (r *AssignClassroomRepo) GetCourse(courseID string) (course.Course, error) {
	for _, item := range r.data.Courses() {
		if item.ID() == courseID {
			return *item, nil
		}
	}
	return course.Course{}, fmt.Errorf("%w: %s", repo.ErrCourseNotFound, courseID)
}
