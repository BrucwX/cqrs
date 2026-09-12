package model

import "time"

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
