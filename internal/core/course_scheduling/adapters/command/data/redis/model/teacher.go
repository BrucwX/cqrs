package model

import (
	"fmt"
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
)

// Teacher 是 Redis 快照（对应 teacher），领域模型是 teacher.Teacher。
//
// 值对象 ContactInfo 就地打平为 Phone / Email。
//
// StudentID 是这位讲师「作为学员上课时」的 ID：报名（CourseEnrollment）与缺勤
// （AbsenceRecord）记录都挂在它上面，授证前「先以学员身份修完课程」也是按它查。
type Teacher struct {
	ID        int64  `json:"id"`         // 讲师 ID
	StudentID int64  `json:"student_id"` // 作为学员上课时的 ID（报名、缺勤记录用）
	Name      string `json:"name"`       // 姓名
	Title     string `json:"title"`      // 职衔（金牌讲师 / 特级培训师）

	Phone string `json:"phone"` // 手机号      <- ContactInfo.Phone
	Email string `json:"email"` // 邮箱（可空）<- ContactInfo.Email

	// 1 在职 / 2 休假 / 3 已离职
	Status uint8 `json:"status"`

	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// TeacherToRedis 写路径：ContactInfo 值对象打平成 Phone / Email 两个字段。
func TeacherToRedis(do *teacher.Teacher) (*Teacher, error) {
	if do == nil {
		return nil, ErrTeacherToRedis
	}
	contact := do.Contact()
	return &Teacher{
		ID:        do.ID(),
		StudentID: do.StudentID(),
		Name:      do.Name(),
		Title:     do.Title(),
		Phone:     contact.Phone(),
		Email:     contact.Email(),
		Status:    uint8(do.Status()),
		CreatedAt: do.CreatedAt(),
		UpdatedAt: do.UpdatedAt(),
	}, nil
}

// TeacherFromRedis 读路径：手机号不符合格式会被值对象构造函数拒绝。
func TeacherFromRedis(po *Teacher) (*teacher.Teacher, error) {
	if po == nil {
		return nil, ErrTeacherFromRedis
	}
	contact, err := teacher.NewContactInfo(po.Phone, po.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTeacherFromRedis, err)
	}
	return teacher.Reconstitute(
		po.ID, po.StudentID, po.Name, po.Title, contact,
		teacher.Status(po.Status), po.CreatedAt, po.UpdatedAt,
	), nil
}
