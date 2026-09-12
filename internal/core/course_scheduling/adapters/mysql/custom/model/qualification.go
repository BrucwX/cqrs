package model

import "time"

// Qualification 对应表 qualification，领域模型是 qualification.Qualification。
//
// 资质绑「课程类型」而不是具体课程：讲师能不能教某门课，等于顺着
// course -> course_type -> qualification 查这条链路。
//
// 能不能发由领域规则判定（讲师得先以学员身份修完该类型下的课），
// 所以表里只有 CertifiedAt 与 Status，没有审核流字段。
//
// Status = 3（已过期）在当前模型里不会出现：领域模型已去掉有效期，
// 取值仍列出来是为了和建表脚本的列注释对齐。
//
// 这张表只有 updated_at，没有 created_at —— 认证时间本身就是 CertifiedAt。
type Qualification struct {
	ID           int64     // bigint      资质 ID
	TeacherID    int64     // bigint      讲师 ID
	CourseTypeID string    // varchar(36) 课程类型 ID（course_type.id）
	CertifiedAt  time.Time // datetime    认证/试讲通过时间

	// uint8 tinyint unsigned 1 生效 / 2 已吊销 / 3 已过期（当前不会出现）
	Status uint8

	UpdatedAt   time.Time // datetime 更新时间
	LockVersion uint64    // bigint unsigned 乐观锁版本号
}
