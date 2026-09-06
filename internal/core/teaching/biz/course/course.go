package course

import (
	"time"

	"github.com/google/uuid"
)

// Course represents a teaching course. It is the course aggregate root.
type Course struct {
	ID              string
	Name            string
	Description     string
	TeacherID       string
	ClassroomID     string
	Semester        Semester       // 学期
	CourseStartTime time.Time      // 课程开始日期，必须与 WeeklySlots 中的某一天匹配
	TotalSessions   int32          // 总课时数
	RemainingSessions int32        // 剩余课时
	WeeklySlots     []ScheduleTime // 每周上课时间段
	MaxStudents     int32
	StudentIDs      map[string]struct{} // 已报名学员
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

// NewCourse creates a new Course aggregate with the given parameters.
// It generates a UUID, initializes internal state, and validates the course.
// courseStartTime must match one of the weekdays in weeklySlots.
func NewCourse(
	name, description string,
	teacherID, classroomID string,
	maxStudents, totalSessions int32,
	semester Semester,
	courseStartTime time.Time,
	weeklySlots []ScheduleTime,
) (*Course, error) {
	now := time.Now()
	c := &Course{
		ID:              uuid.New().String(),
		Name:            name,
		Description:     description,
		TeacherID:       teacherID,
		ClassroomID:     classroomID,
		Semester:        semester,
		CourseStartTime: courseStartTime,
		TotalSessions:   totalSessions,
		RemainingSessions: totalSessions,
		WeeklySlots:     weeklySlots,
		MaxStudents:     maxStudents,
		StudentIDs:      make(map[string]struct{}),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if c.WeeklySlots == nil {
		c.WeeklySlots = make([]ScheduleTime, 0)
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return c, nil
}

// Validate validates the course.
func (c *Course) Validate() error {
	if c.Name == "" {
		return ErrCourseInvalidArgument
	}
	if c.TeacherID == "" {
		return ErrCourseTeacherNotFound
	}
	if c.ClassroomID == "" {
		return ErrCourseClassroomNotFound
	}
	if c.MaxStudents <= 0 {
		return ErrCourseInvalidArgument
	}
	if c.TotalSessions <= 0 {
		return ErrCourseInvalidArgument
	}
	// 学期校验
	if c.Semester.StartTime.IsZero() || c.Semester.EndTime.IsZero() {
		return ErrCourseInvalidArgument
	}
	// 课程开始时间校验
	if c.CourseStartTime.IsZero() {
		return ErrCourseInvalidArgument
	}
	// 课程开始时间必须在学期内
	if !c.Semester.Contains(c.CourseStartTime) {
		return ErrCourseInvalidArgument
	}
	// 课程开始时间必须与 WeeklySlots 中的某一天匹配
	if len(c.WeeklySlots) > 0 {
		courseWeekday := Weekday(c.CourseStartTime.Weekday())
		matched := false
		for _, slot := range c.WeeklySlots {
			if slot.Day == courseWeekday {
				matched = true
				break
			}
		}
		if !matched {
			return ErrCourseInvalidSchedule
		}
	}
	return nil
}

// --- 学员管理 ---

// EnrollStudent enrolls a student in the course.
func (c *Course) EnrollStudent(studentID string) error {
	if c.MaxStudents > 0 && int32(len(c.StudentIDs)) >= c.MaxStudents {
		return ErrCourseFull
	}
	if c.IsEnrolled(studentID) {
		return ErrCourseStudentAlreadyEnrolled
	}
	c.StudentIDs[studentID] = struct{}{}
	return nil
}

// UnenrollStudent unenrolls a student from the course.
func (c *Course) UnenrollStudent(studentID string) error {
	if !c.IsEnrolled(studentID) {
		return ErrCourseStudentNotEnrolled
	}
	delete(c.StudentIDs, studentID)
	return nil
}

// IsEnrolled checks if a student is enrolled in the course.
func (c *Course) IsEnrolled(studentID string) bool {
	_, ok := c.StudentIDs[studentID]
	return ok
}

// StudentCount returns the number of enrolled students.
func (c *Course) StudentCount() int32 {
	return int32(len(c.StudentIDs))
}

// --- 课时管理 ---

// ConsumeSession consumes one session. Called when a class is actually held.
func (c *Course) ConsumeSession() error {
	if c.RemainingSessions <= 0 {
		return ErrCourseNoSessionsLeft
	}
	c.RemainingSessions--
	c.UpdatedAt = time.Now()
	return nil
}

// IsFinished checks if all sessions have been consumed.
func (c *Course) IsFinished() bool {
	return c.RemainingSessions <= 0
}

// --- 每周时间段管理 ---

// AddWeeklySlot adds a weekly time slot for the course.
func (c *Course) AddWeeklySlot(slot ScheduleTime) {
	c.WeeklySlots = append(c.WeeklySlots, slot)
	c.UpdatedAt = time.Now()
}

// RemoveWeeklySlot removes a weekly time slot by matching day and time.
func (c *Course) RemoveWeeklySlot(slot ScheduleTime) {
	filtered := make([]ScheduleTime, 0, len(c.WeeklySlots))
	for _, s := range c.WeeklySlots {
		if !s.Equals(&slot) {
			filtered = append(filtered, s)
		}
	}
	c.WeeklySlots = filtered
	c.UpdatedAt = time.Now()
}

// --- 冲突检测 ---

// HasScheduleConflict checks if this course has a time conflict with another course.
// It checks both semester date range overlap and weekly slot overlap.
func (c *Course) HasScheduleConflict(other *Course) bool {
	if other == nil {
		return false
	}
	// 学期日期范围不重叠则无冲突
	if !c.Semester.Overlaps(&other.Semester) {
		return false
	}
	// 检查每周时间段是否有重叠
	for _, a := range c.WeeklySlots {
		for _, b := range other.WeeklySlots {
			if a.Overlaps(&b) {
				return true
			}
		}
	}
	return false
}

// --- 属性变更 ---

// ChangeTeacher changes the teacher of the course.
func (c *Course) ChangeTeacher(teacherID string) error {
	if teacherID == "" {
		return ErrCourseInvalidArgument
	}
	c.TeacherID = teacherID
	c.UpdatedAt = time.Now()
	return nil
}

// ChangeClassroom changes the classroom of the course.
func (c *Course) ChangeClassroom(classroomID string) error {
	if classroomID == "" {
		return ErrCourseInvalidArgument
	}
	c.ClassroomID = classroomID
	c.UpdatedAt = time.Now()
	return nil
}

// ChangeSemester changes the semester period.
func (c *Course) ChangeSemester(semester Semester) error {
	if semester.StartTime.IsZero() || semester.EndTime.IsZero() {
		return ErrCourseInvalidArgument
	}
	c.Semester = semester
	c.UpdatedAt = time.Now()
	return nil
}

// ChangeCourseStartTime changes the course start date.
// The new start time must match one of the weekdays in WeeklySlots.
func (c *Course) ChangeCourseStartTime(startTime time.Time) error {
	if startTime.IsZero() {
		return ErrCourseInvalidArgument
	}
	// 课程开始时间必须在学期内
	if !c.Semester.Contains(startTime) {
		return ErrCourseInvalidArgument
	}
	// 课程开始时间必须与 WeeklySlots 中的某一天匹配
	if len(c.WeeklySlots) > 0 {
		courseWeekday := Weekday(startTime.Weekday())
		matched := false
		for _, slot := range c.WeeklySlots {
			if slot.Day == courseWeekday {
				matched = true
				break
			}
		}
		if !matched {
			return ErrCourseInvalidSchedule
		}
	}
	c.CourseStartTime = startTime
	c.UpdatedAt = time.Now()
	return nil
}
