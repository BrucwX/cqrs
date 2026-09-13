package model

import (
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
)

// Qualification 是 Redis 快照（对应 qualification），领域模型是 qualification.Qualification。
//
// 资质绑「课程类型」而不是具体课程：讲师能不能教某门课，等于顺着
// course -> course_type -> qualification 查这条链路。
//
// 能不能发由领域规则判定（讲师得先以学员身份修完该类型下的课），
// 所以只有 CertifiedAt 与 Status，没有审核流字段。
//
// Status = 3（已过期）在当前模型里不会出现：领域模型已去掉有效期，
// 取值仍列出来是为了和领域模型对齐。
//
// 快照里只有 updated_at，没有 created_at —— 认证时间本身就是 CertifiedAt。
type Qualification struct {
	ID           int64     `json:"id"`             // 资质 ID
	TeacherID    int64     `json:"teacher_id"`     // 讲师 ID
	CourseTypeID string    `json:"course_type_id"` // 课程类型 ID（course_type.id）
	CertifiedAt  time.Time `json:"certified_at"`   // 认证/试讲通过时间

	// 1 生效 / 2 已吊销 / 3 已过期（当前不会出现）
	Status uint8 `json:"status"`

	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// QualificationToRedis 写路径。快照里只有 updated_at，认证时间就是 CertifiedAt。
func QualificationToRedis(do *qualification.Qualification) (*Qualification, error) {
	if do == nil {
		return nil, ErrQualificationToRedis
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

// QualificationFromRedis 读路径。
func QualificationFromRedis(po *Qualification) (*qualification.Qualification, error) {
	if po == nil {
		return nil, ErrQualificationFromRedis
	}
	return qualification.Reconstitute(
		po.ID, po.TeacherID, po.CourseTypeID, po.CertifiedAt,
		qualification.Status(po.Status), po.UpdatedAt,
	), nil
}
