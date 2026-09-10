package teacher

import (
	"errors"
	"time"
)

// --- 聚合根 (Aggregate Root) ---

type Teacher struct {
	id        int64
	name      string
	title     string      // 职衔（如 "特级培训师"、"金牌讲师"）
	contact   ContactInfo // 联系方式
	status    Status      // 任职状态
	createdAt time.Time
	updatedAt time.Time
}

// NewTeacher 录入新全职讲师（初始化默认为正常在职状态）
func NewTeacher(
	id int64,
	name string,
	title string,
	contact ContactInfo,
) (*Teacher, error) {
	if name == "" {
		return nil, errors.New("teacher name is required")
	}

	now := time.Now()
	return &Teacher{
		id:        id,
		name:      name,
		title:     title,
		contact:   contact,
		status:    StatusActive,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// Reconstitute 从仓储层/数据库还原聚合根
func Reconstitute(
	id int64,
	name string,
	title string,
	contact ContactInfo,
	status Status,
	createdAt, updatedAt time.Time,
) *Teacher {
	return &Teacher{
		id:        id,
		name:      name,
		title:     title,
		contact:   contact,
		status:    status,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// --- 核心领域行为 (Domain Behaviors) ---

// CanAssignSchedule 校验讲师当前是否具备排课准入资格
func (t *Teacher) CanAssignSchedule() error {
	switch t.status {
	case StatusActive:
		return nil
	case StatusOnLeave:
		return ErrTeacherOnLeave
	case StatusTerminated:
		return ErrTeacherNotActive
	default:
		return errors.New("unknown teacher status")
	}
}

// TakeLeave 登记休假（休假期间不可排课）
func (t *Teacher) TakeLeave() error {
	if t.status == StatusTerminated {
		return ErrTeacherNotActive
	}
	if t.status == StatusOnLeave {
		return ErrAlreadyOnLeave
	}
	t.status = StatusOnLeave
	t.updatedAt = time.Now()
	return nil
}

// ResumeWork 销假 / 恢复正常工作
func (t *Teacher) ResumeWork() error {
	if t.status == StatusTerminated {
		return ErrTeacherNotActive
	}
	if t.status == StatusActive {
		return ErrAlreadyActive
	}
	t.status = StatusActive
	t.updatedAt = time.Now()
	return nil
}

// Terminate 离职处理
func (t *Teacher) Terminate() error {
	if t.status == StatusTerminated {
		return ErrAlreadyTerminated
	}
	t.status = StatusTerminated
	t.updatedAt = time.Now()
	return nil
}

// UpdateProfile 更新基础资料与联系方式
func (t *Teacher) UpdateProfile(title string, contact ContactInfo) {
	t.title = title
	t.contact = contact
	t.updatedAt = time.Now()
}

// --- 只读属性访问器 (Getters) ---

func (t *Teacher) ID() int64            { return t.id }
func (t *Teacher) Name() string         { return t.name }
func (t *Teacher) Title() string        { return t.title }
func (t *Teacher) Contact() ContactInfo { return t.contact }
func (t *Teacher) Status() Status       { return t.status }
func (t *Teacher) CreatedAt() time.Time { return t.createdAt }
func (t *Teacher) UpdatedAt() time.Time { return t.updatedAt }

func (t *Teacher) IsActive() bool     { return t.status == StatusActive }
func (t *Teacher) IsOnLeave() bool    { return t.status == StatusOnLeave }
func (t *Teacher) IsTerminated() bool { return t.status == StatusTerminated }
