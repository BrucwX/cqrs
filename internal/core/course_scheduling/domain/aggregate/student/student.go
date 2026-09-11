package student

import (
	"errors"
	"time"
)

// --- 聚合根 (Aggregate Root) ---

type Student struct {
	id          int64 // 全局统一的人员/用户 ID (与 Teacher.ID 保持一致)
	name        string
	studentType StudentType // 内部员工 / 外部客户
	contact     ContactInfo
	status      Status
	createdAt   time.Time
	updatedAt   time.Time
}

// NewStudent 录入学员档案
func NewStudent(
	name string,
	studentType StudentType,
	contact ContactInfo,
) (*Student, error) {
	if name == "" {
		return nil, errors.New("student name is required")
	}

	now := time.Now()
	return &Student{
		id:          generateID(),
		name:        name,
		studentType: studentType,
		contact:     contact,
		status:      StatusActive,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// generateID 生成唯一的 int64 ID
func generateID() int64 {
	return time.Now().UnixNano()
}

// Reconstitute 仓储恢复聚合根
func Reconstitute(
	id int64,
	name string,
	studentType StudentType,
	contact ContactInfo,
	status Status,
	createdAt, updatedAt time.Time,
) *Student {
	return &Student{
		id:          id,
		name:        name,
		studentType: studentType,
		contact:     contact,
		status:      status,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// --- 核心领域行为 (Domain Behaviors) ---

// CanEnrollCourse 校验是否有选课/参训资格，同时防止"既是讲师又当本课学生"
func (s *Student) CanEnrollCourse(courseInstructorID int64) error {
	switch s.status {
	case StatusActive:
		// 规则：讲师不能选报自己主讲的课程
		if s.id == courseInstructorID {
			return ErrSelfEnrollmentProhibited
		}
		return nil
	case StatusBlocked:
		return ErrStudentBlocked
	case StatusCancelled:
		return ErrStudentCancelled
	default:
		return ErrUnknownStatus
	}
}

// IsInternal 是否为公司内部员工
func (s *Student) IsInternal() bool {
	return s.studentType == TypeInternal
}

// UpdateContact 变更联系方式
func (s *Student) UpdateContact(contact ContactInfo) {
	s.contact = contact
	s.updatedAt = time.Now()
}

// Block 封禁 / 暂停学习资格
func (s *Student) Block() {
	s.status = StatusBlocked
	s.updatedAt = time.Now()
}

// Unblock 解封
func (s *Student) Unblock() {
	s.status = StatusActive
	s.updatedAt = time.Now()
}

// Cancel 注销 / 离职结案
func (s *Student) Cancel() {
	s.status = StatusCancelled
	s.updatedAt = time.Now()
}

// --- 只读属性访问器 (Getters) ---

func (s *Student) ID() int64                { return s.id }
func (s *Student) Name() string             { return s.name }
func (s *Student) StudentType() StudentType { return s.studentType }
func (s *Student) Contact() ContactInfo     { return s.contact }
func (s *Student) Status() Status           { return s.status }
func (s *Student) CreatedAt() time.Time     { return s.createdAt }
func (s *Student) UpdatedAt() time.Time     { return s.updatedAt }
