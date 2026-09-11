package command

import (
	"context"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
	"cqrs/internal/core/course_scheduling/domain/aggregate/course"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// CourseSlotCommand 课表槽位命令实现（内存版）。
type CourseSlotCommand struct {
	data *memory.Data
}

// 编译期断言：内存实现必须满足命令接口。
var _ repo.CourseSlotCommand = (*CourseSlotCommand)(nil)

// NewCourseSlotCommand 创建课表槽位命令实现
//
// 返回接口类型，这样实现不完整时会在编译期暴露。
func NewCourseSlotCommand(d *memory.Data) repo.CourseSlotCommand {
	return &CourseSlotCommand{data: d}
}

// Save 保存课表槽位（新增或更新）
func (c *CourseSlotCommand) Save(cs *courseSlot.CourseSlot) error {
	if cs == nil {
		return repo.ErrCourseSlotRequired
	}
	c.data.SaveCourseSlot(cs)
	return nil
}

// Delete 删除课表槽位
func (c *CourseSlotCommand) Delete(id string) error {
	if !c.data.DeleteCourseSlot(id) {
		return fmt.Errorf("%w: %s", repo.ErrCourseSlotNotFound, id)
	}
	return nil
}

// AssignTeacher 给指定课表槽位们配置老师
//
// 先把该讲师聚合取出来交给调用方判定；传 nil 表示不做冲突检查。
// 冲突时整批中止，不写入任何槽位。
func (c *CourseSlotCommand) AssignTeacher(ctx context.Context, slotIDs []string, teacherID int64, checkConflictFn func(ctx context.Context, slots []courseSlot.CourseSlot, t teacher.Teacher) (bool, error)) error {
	slots, err := c.loadSlots(slotIDs)
	if err != nil {
		return err
	}

	if checkConflictFn != nil {
		teacher, err := c.teacherByID(teacherID)
		if err != nil {
			return err
		}

		conflict, err := checkConflictFn(ctx, slotValues(slots), teacher)
		if err != nil {
			return err
		}
		if conflict {
			return fmt.Errorf("%w: teacher %d", repo.ErrCourseSlotConflict, teacherID)
		}
	}

	for _, cs := range slots {
		if err := cs.ChangeTeacher(teacherID); err != nil {
			return err
		}
		c.data.SaveCourseSlot(cs)
	}
	return nil
}

// AssignCourse 给指定课表槽位们配置课程
//
// 先把课程聚合取出来交给调用方判定；传 nil 表示不做冲突检查。
// 冲突时整批中止，不写入任何槽位。
func (c *CourseSlotCommand) AssignCourse(ctx context.Context, slotIDs []string, courseID string, checkConflictFn func(ctx context.Context, slots []courseSlot.CourseSlot, crs course.Course) (bool, error)) error {
	slots, err := c.loadSlots(slotIDs)
	if err != nil {
		return err
	}

	if checkConflictFn != nil {
		crs, err := c.courseByID(courseID)
		if err != nil {
			return err
		}

		conflict, err := checkConflictFn(ctx, slotValues(slots), crs)
		if err != nil {
			return err
		}
		if conflict {
			return fmt.Errorf("%w: course %s", repo.ErrCourseSlotConflict, courseID)
		}
	}

	for _, cs := range slots {
		if err := cs.ChangeCourse(courseID); err != nil {
			return err
		}
		c.data.SaveCourseSlot(cs)
	}
	return nil
}

// AssignClassroom 给指定课表槽位们配置教室
//
// 先把教室聚合取出来交给调用方判定；传 nil 表示不做冲突检查。
// 冲突时整批中止，不写入任何槽位。
func (c *CourseSlotCommand) AssignClassroom(ctx context.Context, slotIDs []string, classroomID string, checkConflictFn func(ctx context.Context, slots []courseSlot.CourseSlot, cr classroom.Classroom) (bool, error)) error {
	slots, err := c.loadSlots(slotIDs)
	if err != nil {
		return err
	}

	if checkConflictFn != nil {
		cr, err := c.classroomByID(classroomID)
		if err != nil {
			return err
		}

		conflict, err := checkConflictFn(ctx, slotValues(slots), cr)
		if err != nil {
			return err
		}
		if conflict {
			return fmt.Errorf("%w: classroom %s", repo.ErrCourseSlotConflict, classroomID)
		}
	}

	for _, cs := range slots {
		if err := cs.ChangeClassroom(classroomID); err != nil {
			return err
		}
		c.data.SaveCourseSlot(cs)
	}
	return nil
}

// --- 内部实现 ---

// teacherByID 取讲师聚合，取不到报 not found。
func (c *CourseSlotCommand) teacherByID(id int64) (teacher.Teacher, error) {
	for _, item := range c.data.Teachers() {
		if item.ID() == id {
			return *item, nil
		}
	}
	return teacher.Teacher{}, fmt.Errorf("%w: %d", repo.ErrTeacherNotFound, id)
}

// loadSlots 按 ID 取出槽位指针；任何一个不存在就整体失败（不做部分写入）。
func (c *CourseSlotCommand) loadSlots(slotIDs []string) ([]*courseSlot.CourseSlot, error) {
	slots := make([]*courseSlot.CourseSlot, 0, len(slotIDs))
	for _, id := range slotIDs {
		cs, ok := c.data.CourseSlotByID(id)
		if !ok {
			return nil, fmt.Errorf("%w: %s", repo.ErrCourseSlotNotFound, id)
		}
		slots = append(slots, cs)
	}
	return slots, nil
}

// slotValues 把槽位指针切片转成值拷贝切片，避免调用方误改聚合。
func slotValues(slots []*courseSlot.CourseSlot) []courseSlot.CourseSlot {
	out := make([]courseSlot.CourseSlot, 0, len(slots))
	for _, cs := range slots {
		out = append(out, *cs)
	}
	return out
}

// classroomByID 取教室聚合，取不到报 not found。
func (c *CourseSlotCommand) classroomByID(id string) (classroom.Classroom, error) {
	for _, item := range c.data.Classrooms() {
		if item.ID() == id {
			return *item, nil
		}
	}
	return classroom.Classroom{}, fmt.Errorf("%w: %s", repo.ErrClassroomNotFound, id)
}

// courseByID 取课程聚合，取不到报 not found。
func (c *CourseSlotCommand) courseByID(id string) (course.Course, error) {
	for _, item := range c.data.Courses() {
		if item.ID() == id {
			return *item, nil
		}
	}
	return course.Course{}, fmt.Errorf("%w: %s", repo.ErrCourseNotFound, id)
}

// indexBy 按 key 建索引；重复 key 保留先出现的一个。
func indexBy[K comparable, V any](items []V, key func(V) K) map[K]V {
	out := make(map[K]V, len(items))
	for _, item := range items {
		k := key(item)
		if _, ok := out[k]; !ok {
			out[k] = item
		}
	}
	return out
}
