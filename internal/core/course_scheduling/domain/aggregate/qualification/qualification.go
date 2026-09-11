package qualification

import (
	"errors"
	"time"
)

// --- 聚合根 (Aggregate Root) ---

type Qualification struct {
	id           int64
	teacherID    int64     // 关联的讲师 ID
	courseTypeID string    // 关联的课程类型 ID（courseType.CourseType）
	certifiedAt  time.Time // 认证/试讲通过时间
	expireAt     time.Time // 资质有效期截止时间（若长期有效可设为零值或未来远期时间）
	status       Status    // 资质状态
	updatedAt    time.Time
}

// generateID 生成唯一的 int64 ID
func generateID() int64 {
	return time.Now().UnixNano()
}

// NewQualification 颁发/授予授课资质
//
// 资质绑定的是「课程类型」而不是具体课程：讲师通过某类课程的试讲，
// 即可讲授该类型下的所有课程。
func NewQualification(
	teacherID int64,
	courseTypeID string,
	certifiedAt time.Time,
	expireAt time.Time,
) (*Qualification, error) {
	if teacherID <= 0 {
		return nil, errors.New("invalid teacher ID")
	}
	if courseTypeID == "" {
		return nil, errors.New("course type ID is required")
	}

	return &Qualification{
		id:           generateID(),
		teacherID:    teacherID,
		courseTypeID: courseTypeID,
		certifiedAt:  certifiedAt,
		expireAt:     expireAt,
		status:       StatusActive,
		updatedAt:    time.Now(),
	}, nil
}

// Reconstitute 从仓储层还原聚合根
func Reconstitute(
	id int64,
	teacherID int64,
	courseTypeID string,
	certifiedAt time.Time,
	expireAt time.Time,
	status Status,
	updatedAt time.Time,
) *Qualification {
	return &Qualification{
		id:           id,
		teacherID:    teacherID,
		courseTypeID: courseTypeID,
		certifiedAt:  certifiedAt,
		expireAt:     expireAt,
		status:       status,
		updatedAt:    updatedAt,
	}
}

// --- 核心领域行为 (Domain Behaviors) ---

// IsEligible 检查该老师在指定时间点是否有资格讲授该门课
func (q *Qualification) IsEligible(at time.Time) error {
	if q.status == StatusRevoked {
		return ErrQualificationRevoked
	}
	// 如果设置了过期时间，校验是否已超时
	if !q.expireAt.IsZero() && at.After(q.expireAt) {
		return ErrQualificationExpired
	}
	return nil
}

// Revoke 吊销授课资格（例如试讲不达标或客诉严重取消授课授权）
func (q *Qualification) Revoke() {
	q.status = StatusRevoked
	q.updatedAt = time.Now()
}

// Renew 资质续期
func (q *Qualification) Renew(newExpireAt time.Time) {
	q.expireAt = newExpireAt
	q.status = StatusActive
	q.updatedAt = time.Now()
}

// --- 只读属性访问器 (Getters) ---

func (q *Qualification) ID() int64              { return q.id }
func (q *Qualification) TeacherID() int64       { return q.teacherID }
func (q *Qualification) CourseTypeID() string   { return q.courseTypeID }
func (q *Qualification) CertifiedAt() time.Time { return q.certifiedAt }
func (q *Qualification) ExpireAt() time.Time    { return q.expireAt }
func (q *Qualification) Status() Status         { return q.status }
