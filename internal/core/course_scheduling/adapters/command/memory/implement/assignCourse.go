package implement

import (
	"cqrs/internal/core/course_scheduling/adapters/command/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// AssignCourseRepo 是 repo.AssignCourseRepo 的内存实现。
//
// 它只做「按需取数据」：配课程的规则判定在命令服务里（checkCourse），
// 这里不参与任何规则。
type AssignCourseRepo struct {
	data *memory.Data
}

// 编译期断言：内存实现必须满足接口。
var _ repo.AssignCourseRepo = (*AssignCourseRepo)(nil)

// NewAssignCourseRepo 创建配课程冲突检查的数据来源。
func NewAssignCourseRepo(d *memory.Data) repo.AssignCourseRepo {
	return &AssignCourseRepo{data: d}
}

// GetCourseSlots 取该课程现有的全部排期。
func (r *AssignCourseRepo) GetCourseSlots(courseID string) (courseSlot.CourseSlots, error) {
	out := make(courseSlot.CourseSlots, 0)
	for _, cs := range r.data.CourseSlots() {
		if cs.CourseID() == courseID {
			out = append(out, *cs)
		}
	}
	return out, nil
}
