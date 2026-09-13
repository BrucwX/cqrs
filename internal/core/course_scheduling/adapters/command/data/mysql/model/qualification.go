package model

import (
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
)

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

// QualificationToPO 写路径。表里只有 updated_at，认证时间就是 CertifiedAt。
func QualificationToPO(do *qualification.Qualification) (*Qualification, error) {
	if do == nil {
		return nil, ErrQualificationDOToPO
	}
	return &Qualification{
		ID:           do.ID(),
		TeacherID:    do.TeacherID(),
		CourseTypeID: do.CourseTypeID(),
		CertifiedAt:  do.CertifiedAt(),
		Status:       uint8(do.Status()),
		UpdatedAt:    do.UpdatedAt(),
	}, nil
}

// QualificationToDO 读路径。
func QualificationToDO(po *Qualification) (*qualification.Qualification, error) {
	if po == nil {
		return nil, ErrQualificationPOToDO
	}
	return qualification.Reconstitute(
		po.ID, po.TeacherID, po.CourseTypeID, po.CertifiedAt,
		qualification.Status(po.Status), po.UpdatedAt,
	), nil
}
