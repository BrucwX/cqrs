package implement

import (
	"context"
	"fmt"

	"cqrs/internal/core/course_scheduling/adapters/command/memory"
	"cqrs/internal/core/course_scheduling/domain/aggregate/courseSlot"
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
func (c *CourseSlotCommand) Save(ctx context.Context, cs *courseSlot.CourseSlot) error {
	if cs == nil {
		return courseSlot.ErrCourseSlotRequired
	}
	c.data.SaveCourseSlot(cs)
	return nil
}

// Delete 删除课表槽位
func (c *CourseSlotCommand) Delete(ctx context.Context, id string) error {
	if !c.data.DeleteCourseSlot(id) {
		return fmt.Errorf("%w: %s", courseSlot.ErrCourseSlotNotFound, id)
	}
	return nil
}

// MustGet 取课表槽位聚合；不存在时报 ErrCourseSlotNotFound。
func (c *CourseSlotCommand) MustGet(ctx context.Context, id string) (courseSlot.CourseSlot, error) {
	item, ok := c.data.CourseSlotByID(id)
	if !ok {
		return courseSlot.CourseSlot{}, fmt.Errorf("%w: %s", courseSlot.ErrCourseSlotNotFound, id)
	}
	return *item, nil
}

// AssignTeacher 给指定课表槽位们配置老师
//
// 只写。判定由调用方（领域服务）做完才调过来，仓库不认识规则。
func (c *CourseSlotCommand) AssignTeacher(ctx context.Context, slotIDs []string, teacherID int64) error {
	slots, err := c.loadSlots(slotIDs)
	if err != nil {
		return err
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
// 只写。判定由调用方（领域服务）做完才调过来。
func (c *CourseSlotCommand) AssignCourse(ctx context.Context, slotIDs []string, courseID string) error {
	slots, err := c.loadSlots(slotIDs)
	if err != nil {
		return err
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
// 只写。判定由调用方（领域服务）做完才调过来。
func (c *CourseSlotCommand) AssignClassroom(ctx context.Context, slotIDs []string, classroomID string) error {
	slots, err := c.loadSlots(slotIDs)
	if err != nil {
		return err
	}

	for _, cs := range slots {
		if err := cs.ChangeClassroom(classroomID); err != nil {
			return err
		}
		c.data.SaveCourseSlot(cs)
	}
	return nil
}

// --- 排期读取（按返回值的聚合根归到本接口）---

// GetSlots 按 ID 取本次要排的槽位，顺序与入参一致；少一个就报 not found。
func (c *CourseSlotCommand) GetSlots(ctx context.Context, slotIDs []string) (courseSlot.CourseSlots, error) {
	out := make(courseSlot.CourseSlots, 0, len(slotIDs))
	for _, id := range slotIDs {
		item, ok := c.data.CourseSlotByID(id)
		if !ok {
			return nil, fmt.Errorf("%w: %s", courseSlot.ErrCourseSlotNotFound, id)
		}
		out = append(out, *item)
	}
	return out, nil
}

// GetTeacherSlots 取该讲师现有的全部排期。
func (c *CourseSlotCommand) GetTeacherSlots(ctx context.Context, teacherID int64) (courseSlot.CourseSlots, error) {
	out := make(courseSlot.CourseSlots, 0)
	for _, cs := range c.data.CourseSlots() {
		if cs.TeacherID() == teacherID {
			out = append(out, *cs)
		}
	}
	return out, nil
}

// GetClassroomSlots 取该教室现有的全部排期。
func (c *CourseSlotCommand) GetClassroomSlots(ctx context.Context, classroomID string) (courseSlot.CourseSlots, error) {
	out := make(courseSlot.CourseSlots, 0)
	for _, cs := range c.data.CourseSlots() {
		if cs.ClassroomID() == classroomID {
			out = append(out, *cs)
		}
	}
	return out, nil
}

// GetCourseSlots 取该课程现有的全部排期。
func (c *CourseSlotCommand) GetCourseSlots(ctx context.Context, courseID string) (courseSlot.CourseSlots, error) {
	return c.slotsOfCourses(courseID), nil
}

// GetStudentSlots 取该学员现有「在学」课程的全部排期。
//
// 已退课 / 已结业的记录不计入；结果里会包含目标课程自身的排期，
// 所以重复选课会表现为与自身槽位重叠。
func (c *CourseSlotCommand) GetStudentSlots(ctx context.Context, studentID int64) (courseSlot.CourseSlots, error) {
	courseIDs := make([]string, 0)
	for _, e := range c.data.Enrollments() {
		if e.StudentID() == studentID && e.IsActive() {
			courseIDs = append(courseIDs, e.CourseID())
		}
	}
	return c.slotsOfCourses(courseIDs...), nil
}

// slotsOfCourses 取若干门课程的全部排期。
func (c *CourseSlotCommand) slotsOfCourses(courseIDs ...string) courseSlot.CourseSlots {
	want := make(map[string]struct{}, len(courseIDs))
	for _, id := range courseIDs {
		want[id] = struct{}{}
	}

	out := make(courseSlot.CourseSlots, 0)
	for _, cs := range c.data.CourseSlots() {
		if _, ok := want[cs.CourseID()]; ok {
			out = append(out, *cs)
		}
	}
	return out
}

// --- 内部实现 ---

// loadSlots 按 ID 取出槽位指针；任何一个不存在就整体失败（不做部分写入）。
func (c *CourseSlotCommand) loadSlots(slotIDs []string) ([]*courseSlot.CourseSlot, error) {
	slots := make([]*courseSlot.CourseSlot, 0, len(slotIDs))
	for _, id := range slotIDs {
		cs, ok := c.data.CourseSlotByID(id)
		if !ok {
			return nil, fmt.Errorf("%w: %s", courseSlot.ErrCourseSlotNotFound, id)
		}
		slots = append(slots, cs)
	}
	return slots, nil
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
