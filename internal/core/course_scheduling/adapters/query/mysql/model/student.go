package model

import (
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
)

// Student 对应表 student，领域模型是 student.Student。
//
// 值对象 ContactInfo 就地打平为 Phone / Email。
//
// ID 是全局统一的人员 ID（不是单独的学员自增主键）：内部员工学员
// （StudentType = 2）可能同时是某条 Teacher 记录本人，见 Teacher.StudentID。
type Student struct {
	ID   int64  // bigint      学员 ID（人员/用户 ID）
	Name string // varchar(64) 姓名

	// uint8 tinyint unsigned 1 外部客户学员 / 2 内部员工学员
	StudentType uint8

	Phone string // varchar(20)  手机号      <- ContactInfo.Phone
	Email string // varchar(128) 邮箱（可空）<- ContactInfo.Email

	// uint8 tinyint unsigned 1 正常 / 2 封禁 / 3 已注销
	Status uint8

	CreatedAt   time.Time // datetime 创建时间
	UpdatedAt   time.Time // datetime 更新时间
	LockVersion uint64    // bigint unsigned 乐观锁版本号
}

// StudentToPO 写路径：ContactInfo 值对象打平成 Phone / Email 两列。
func StudentToPO(do *student.Student) *Student {
	contact := do.Contact()
	return &Student{
		ID:          do.ID(),
		Name:        do.Name(),
		StudentType: uint8(do.StudentType()),
		Phone:       contact.Phone(),
		Email:       contact.Email(),
		Status:      uint8(do.Status()),
		CreatedAt:   do.CreatedAt(),
		UpdatedAt:   do.UpdatedAt(),
	}
}

// StudentToDO 读路径：手机号不符合格式会被值对象构造函数拒绝。
func StudentToDO(po *Student) (*student.Student, error) {
	contact, err := student.NewContactInfo(po.Phone, po.Email)
	if err != nil {
		return nil, err
	}
	return student.Reconstitute(
		po.ID, po.Name, student.StudentType(po.StudentType), contact,
		student.Status(po.Status), po.CreatedAt, po.UpdatedAt,
	), nil
}
