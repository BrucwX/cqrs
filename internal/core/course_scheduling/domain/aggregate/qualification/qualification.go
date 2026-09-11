package qualification

import (
	"time"
)

// --- 聚合根 (Aggregate Root) ---

type Qualification struct {
	id           int64
	teacherID    int64     // 关联的讲师 ID
	courseTypeID string    // 关联的课程类型 ID（courseType.CourseType）
	certifiedAt  time.Time // 认证/试讲通过时间
	status       Status    // 资质状态
	updatedAt    time.Time
}

// generateID 生成唯一的 int64 ID
func generateID() int64 {
	return time.Now().UnixNano()
}

// NewQualification 颁发/授予授课资质（ID 由聚合自己生成）
//
// 资质绑定的是「课程类型」而不是具体课程：讲师通过某类课程的试讲，
// 即可讲授该类型下的所有课程。
//
// 讲师够不够格（是否已修完该类型课程）由命令服务判定，聚合只管发证。
func NewQualification(teacherID int64, courseTypeID string) (*Qualification, error) {
	if teacherID <= 0 {
		return nil, ErrTeacherRequired
	}
	if courseTypeID == "" {
		return nil, ErrCourseTypeRequired
	}

	now := time.Now()
	return &Qualification{
		id:           generateID(),
		teacherID:    teacherID,
		courseTypeID: courseTypeID,
		certifiedAt:  now,
		status:       StatusActive,
		updatedAt:    now,
	}, nil
}

// Reconstitute 从仓储层还原聚合根
func Reconstitute(
	id int64,
	teacherID int64,
	courseTypeID string,
	certifiedAt time.Time,
	status Status,
	updatedAt time.Time,
) *Qualification {
	return &Qualification{
		id:           id,
		teacherID:    teacherID,
		courseTypeID: courseTypeID,
		certifiedAt:  certifiedAt,
		status:       status,
		updatedAt:    updatedAt,
	}
}

// --- 核心领域行为 (Domain Behaviors) ---

// IsEligible 检查该老师是否仍有资格讲授该门课（被吊销就没资格了）
func (q *Qualification) IsEligible() error {
	if q.status == StatusRevoked {
		return ErrQualificationRevoked
	}
	return nil
}

// Revoke 吊销授课资格（例如试讲不达标或客诉严重取消授课授权）
func (q *Qualification) Revoke() {
	q.status = StatusRevoked
	q.updatedAt = time.Now()
}

// --- 只读属性访问器 (Getters) ---

func (q *Qualification) ID() int64              { return q.id }
func (q *Qualification) TeacherID() int64       { return q.teacherID }
func (q *Qualification) CourseTypeID() string   { return q.courseTypeID }
func (q *Qualification) CertifiedAt() time.Time { return q.certifiedAt }
func (q *Qualification) Status() Status         { return q.status }
