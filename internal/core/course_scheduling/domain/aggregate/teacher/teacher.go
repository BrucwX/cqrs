package teacher

import (
	"errors"
	"time"
)

// --- 聚合根 (Aggregate Root) ---

type Teacher struct {
	id        int64
	studentID int64 // 讲师作为学员上课时的 ID（报名/缺勤记录挂在它上面）
	name      string
	title     string      // 职衔（如 "特级培训师"、"金牌讲师"）
	contact   ContactInfo // 联系方式
	status    Status      // 任职状态
	createdAt time.Time
	updatedAt time.Time
}

// NewTeacher 录入新讲师（ID 由聚合自己生成）
//
// name 必填；title / contact 不给就取零值；status 不给默认「在职」。
//
// studentID 是讲师作为学员上课时的 ID（报名、缺勤记录都挂在它上面）。
// 学员系统已经给他建过档就传进来，不给（nil 或 <=0）就由聚合自己发一个。
func NewTeacher(
	studentID *int64,
	name *string,
	title *string,
	contact *ContactInfo,
	status *Status,
) (*Teacher, error) {
	if name == nil || *name == "" {
		return nil, errors.New("teacher name is required")
	}

	sid := orZero(studentID)
	if sid <= 0 {
		sid = generateID()
	}

	now := time.Now()
	created := &Teacher{
		id:        generateID(),
		studentID: sid,
		name:      *name,
		title:     orZero(title),
		contact:   orZero(contact),
		status:    StatusActive,
		createdAt: now,
		updatedAt: now,
	}

	if status != nil {
		if err := created.ChangeStatus(*status); err != nil {
			return nil, err
		}
	}
	return created, nil
}

// orZero 解引用可选字段；nil 时取零值。
func orZero[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

// generateID 生成唯一的 int64 ID
func generateID() int64 {
	return time.Now().UnixNano()
}

// Reconstitute 从仓储层/数据库还原聚合根
func Reconstitute(
	id int64,
	studentID int64,
	name string,
	title string,
	contact ContactInfo,
	status Status,
	createdAt, updatedAt time.Time,
) *Teacher {
	return &Teacher{
		id:        id,
		studentID: studentID,
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

// ChangeStatus 切换任职状态
//
// 已经是目标状态时直接返回，保证重复调用幂等；
// 离职后不可再改（Terminate 的守卫仍然生效）。
func (t *Teacher) ChangeStatus(status Status) error {
	if t.status == status {
		return nil
	}

	switch status {
	case StatusActive:
		return t.ResumeWork()
	case StatusOnLeave:
		return t.TakeLeave()
	case StatusTerminated:
		return t.Terminate()
	default:
		return ErrUnknownStatus
	}
}

// Update 按非 nil 的字段更新讲师，nil 的字段保持原值。
//
// 先切状态（唯一可能失败的一步），再改资料，避免失败时留下改了一半的对象。
func (t *Teacher) Update(title *string, contact *ContactInfo, status *Status) error {
	if status != nil {
		if err := t.ChangeStatus(*status); err != nil {
			return err
		}
	}

	if title != nil || contact != nil {
		newTitle := t.title
		if title != nil {
			newTitle = *title
		}
		newContact := t.contact
		if contact != nil {
			newContact = *contact
		}
		t.UpdateProfile(newTitle, newContact)
	}
	return nil
}

// --- 只读属性访问器 (Getters) ---

func (t *Teacher) ID() int64            { return t.id }
func (t *Teacher) StudentID() int64     { return t.studentID }
func (t *Teacher) Name() string         { return t.name }
func (t *Teacher) Title() string        { return t.title }
func (t *Teacher) Contact() ContactInfo { return t.contact }
func (t *Teacher) Status() Status       { return t.status }
func (t *Teacher) CreatedAt() time.Time { return t.createdAt }
func (t *Teacher) UpdatedAt() time.Time { return t.updatedAt }

func (t *Teacher) IsActive() bool     { return t.status == StatusActive }
func (t *Teacher) IsOnLeave() bool    { return t.status == StatusOnLeave }
func (t *Teacher) IsTerminated() bool { return t.status == StatusTerminated }
