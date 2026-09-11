package courseSlot

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// 待分配常量
const (
	// PendingTeacherID 待分配的讲师 ID
	PendingTeacherID int64 = -1
	// PendingClassroomID 待分配的教室 ID
	PendingClassroomID string = ""
)

// --- 聚合根 (Aggregate Root) ---

type CourseSlot struct {
	id          string       // UUID
	courseID    string       // 关联的课程 ID
	weekday     time.Weekday // 星期几（time.Sunday = 0, time.Monday = 1 ... time.Wednesday = 3）
	timeRange   DayTimeRange // 当天的上课时间区间（如 16:00 - 18:00）
	teacherID   int64        // 默认授课讲师 ID（PendingTeacherID 表示待分配）
	classroomID string       // 默认上课教室 ID（PendingClassroomID 表示待分配）
	createdAt   time.Time
	updatedAt   time.Time
}

// NewCourseSlot 创建每周重复的课表排课模板
func NewCourseSlot(
	courseID string,
	weekday time.Weekday,
	timeRange DayTimeRange,
	teacherID int64,
	classroomID string,
) (*CourseSlot, error) {
	if courseID == "" {
		return nil, errors.New("course ID is required")
	}

	now := time.Now()
	return &CourseSlot{
		id:          uuid.New().String(),
		courseID:    courseID,
		weekday:     weekday,
		timeRange:   timeRange,
		teacherID:   teacherID,
		classroomID: classroomID,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// Reconstitute 仓储恢复
func Reconstitute(
	id string,
	courseID string,
	weekday time.Weekday,
	timeRange DayTimeRange,
	teacherID int64,
	classroomID string,
	createdAt, updatedAt time.Time,
) *CourseSlot {
	return &CourseSlot{
		id:          id,
		courseID:    courseID,
		weekday:     weekday,
		timeRange:   timeRange,
		teacherID:   teacherID,
		classroomID: classroomID,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// --- 核心领域行为 (Domain Behaviors) ---

// IsConflictingWith 校验模板层面的冲突（必须星期几相同且日内时间重叠）
func (cs *CourseSlot) IsConflictingWith(other *CourseSlot) bool {
	if cs.weekday != other.weekday {
		return false
	}
	if !cs.timeRange.Overlaps(other.timeRange) {
		return false
	}
	// 时间重合时，同讲师或同教室即为排课模板冲突
	return cs.teacherID == other.teacherID || cs.classroomID == other.classroomID
}

// OverlapsWith 判断两个槽位是否撞时间：星期几相同 且 日内时间段重叠。
//
// 只看时间，不看讲师/教室 —— 需要判"同人/同教室撞车"请用 IsConflictingWith。
func (cs *CourseSlot) OverlapsWith(other CourseSlot) bool {
	return cs.weekday == other.weekday && cs.timeRange.Overlaps(other.timeRange)
}

// UpdateSchedule 调整排课模板的时间与场地
func (cs *CourseSlot) UpdateSchedule(weekday time.Weekday, timeRange DayTimeRange, classroomID string) error {
	if classroomID == "" {
		return errors.New("classroom ID cannot be empty")
	}
	cs.weekday = weekday
	cs.timeRange = timeRange
	cs.classroomID = classroomID
	cs.updatedAt = time.Now()
	return nil
}

// ChangeTeacher 更换该槽位的负责老师
func (cs *CourseSlot) ChangeTeacher(teacherID int64) error {
	if teacherID <= 0 {
		return errors.New("invalid teacher ID")
	}
	cs.teacherID = teacherID
	cs.updatedAt = time.Now()
	return nil
}

// ChangeClassroom 更换该槽位的上课教室
func (cs *CourseSlot) ChangeClassroom(classroomID string) error {
	if classroomID == "" {
		return errors.New("classroom ID cannot be empty")
	}
	cs.classroomID = classroomID
	cs.updatedAt = time.Now()
	return nil
}

// ChangeCourse 更换该槽位所属的课程
func (cs *CourseSlot) ChangeCourse(courseID string) error {
	if courseID == "" {
		return errors.New("course ID cannot be empty")
	}
	cs.courseID = courseID
	cs.updatedAt = time.Now()
	return nil
}

// IsPendingTeacher 检查讲师是否待分配
func (cs *CourseSlot) IsPendingTeacher() bool {
	return cs.teacherID == PendingTeacherID
}

// IsPendingClassroom 检查教室是否待分配
func (cs *CourseSlot) IsPendingClassroom() bool {
	return cs.classroomID == PendingClassroomID
}

// IsFullyAssigned 检查是否已完全分配（讲师和教室都已分配）
func (cs *CourseSlot) IsFullyAssigned() bool {
	return !cs.IsPendingTeacher() && !cs.IsPendingClassroom()
}

// InstantiateForDate 核心行为：将模板按具体的日期实例化为具体的实际时间区间
// 例如：给定某周三的具体日期 2026-09-16，生成对应的 start: 2026-09-16 16:00, end: 2026-09-16 18:00
func (cs *CourseSlot) InstantiateForDate(date time.Time, loc *time.Location) (time.Time, time.Time, error) {
	if date.Weekday() != cs.weekday {
		return time.Time{}, time.Time{}, fmt.Errorf("date %s does not match slot weekday %s", date.Format("2006-01-02"), cs.weekday)
	}

	start := time.Date(
		date.Year(), date.Month(), date.Day(),
		cs.timeRange.start.hour, cs.timeRange.start.minute, 0, 0, loc,
	)
	end := time.Date(
		date.Year(), date.Month(), date.Day(),
		cs.timeRange.end.hour, cs.timeRange.end.minute, 0, 0, loc,
	)

	return start, end, nil
}

// --- 只读属性访问器 (Getters) ---

func (cs *CourseSlot) ID() string              { return cs.id }
func (cs *CourseSlot) CourseID() string        { return cs.courseID }
func (cs *CourseSlot) Weekday() time.Weekday   { return cs.weekday }
func (cs *CourseSlot) TimeRange() DayTimeRange { return cs.timeRange }
func (cs *CourseSlot) TeacherID() int64        { return cs.teacherID }
func (cs *CourseSlot) ClassroomID() string     { return cs.classroomID }
