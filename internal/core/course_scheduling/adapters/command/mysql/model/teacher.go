package model

import (
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
)

// Teacher 对应表 teacher，领域模型是 teacher.Teacher。
//
// 值对象 ContactInfo 就地打平为 Phone / Email。
//
// StudentID 是这位讲师「作为学员上课时」的 ID：报名（CourseEnrollment）与缺勤
// （AbsenceRecord）记录都挂在它上面，授证前「先以学员身份修完课程」也是按它查。
type Teacher struct {
	ID        int64  // bigint      讲师 ID
	StudentID int64  // bigint      作为学员上课时的 ID（报名、缺勤记录用）
	Name      string // varchar(64) 姓名
	Title     string // varchar(64) 职衔（金牌讲师 / 特级培训师）

	Phone string // varchar(20)  手机号      <- ContactInfo.Phone
	Email string // varchar(128) 邮箱（可空）<- ContactInfo.Email

	// uint8 tinyint unsigned 1 在职 / 2 休假 / 3 已离职
	Status uint8

	CreatedAt   time.Time // datetime 创建时间
	UpdatedAt   time.Time // datetime 更新时间
	LockVersion uint64    // bigint unsigned 乐观锁版本号
}

// TeacherToPO 写路径：ContactInfo 值对象打平成 Phone / Email 两列。
func TeacherToPO(do *teacher.Teacher) *Teacher {
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
	}
}

// TeacherToDO 读路径：手机号不符合格式会被值对象构造函数拒绝。
func TeacherToDO(po *Teacher) (*teacher.Teacher, error) {
	contact, err := teacher.NewContactInfo(po.Phone, po.Email)
	if err != nil {
		return nil, err
	}
	return teacher.Reconstitute(
		po.ID, po.StudentID, po.Name, po.Title, contact,
		teacher.Status(po.Status), po.CreatedAt, po.UpdatedAt,
	), nil
}
