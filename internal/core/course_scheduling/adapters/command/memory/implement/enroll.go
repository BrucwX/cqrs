package implement

import (
	"cqrs/internal/core/course_scheduling/adapters/command/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// EnrollRepo 是 repo.EnrollRepo 的内存实现。
//
// 它只做「按需取数据」：选课的规则判定在命令服务里（checkEnrollment），
// 这里不参与任何规则。
type EnrollRepo struct {
	data *memory.Data
}

// 编译期断言：内存实现必须满足接口。
var _ repo.EnrollRepo = (*EnrollRepo)(nil)

// NewEnrollRepo 创建选课冲突检查的数据来源。
func NewEnrollRepo(d *memory.Data) repo.EnrollRepo {
	return &EnrollRepo{data: d}
}

// GetCourseSlots 取某门课程现有的全部排期。
func (r *EnrollRepo) GetCourseSlots(courseID string) (courseSlot.CourseSlots, error) {
	return r.slotsOfCourses(courseID), nil
}

// GetStudentSlots 取该学员现有「在学」课程的全部排期。
//
// 已退课 / 已结业的记录不计入；结果里会包含目标课程自身的排期，
// 所以重复选课会表现为与自身槽位重叠。
func (r *EnrollRepo) GetStudentSlots(studentID int64) (courseSlot.CourseSlots, error) {
	courseIDs := make([]string, 0)
	for _, e := range r.data.Enrollments() {
		if e.StudentID() == studentID && e.IsActive() {
			courseIDs = append(courseIDs, e.CourseID())
		}
	}
	return r.slotsOfCourses(courseIDs...), nil
}

// slotsOfCourses 取若干门课程的全部排期。
func (r *EnrollRepo) slotsOfCourses(courseIDs ...string) courseSlot.CourseSlots {
	want := make(map[string]struct{}, len(courseIDs))
	for _, id := range courseIDs {
		want[id] = struct{}{}
	}

	out := make(courseSlot.CourseSlots, 0)
	for _, cs := range r.data.CourseSlots() {
		if _, ok := want[cs.CourseID()]; ok {
			out = append(out, *cs)
		}
	}
	return out
}
