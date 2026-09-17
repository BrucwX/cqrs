package courseType

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// --- 聚合根 (Aggregate Root) ---

// CourseType 课程类型聚合根。
//
// 它描述「教什么」（如 少儿编程 / 成人英语），而不是「哪个班」。
// 讲师资质绑定的是课程类型（Qualification.CourseTypeID），
// 具体课程（course.Course）通过 CourseTypeID 归属到某个类型 ——
// 这样一位讲师的资质天然覆盖该类型下的所有课程。
type CourseType struct {
	id          string
	name        string
	description string
	createdAt   time.Time
	updatedAt   time.Time
}

// NewCourseType 创建课程类型（ID 由聚合生成）。
func NewCourseType(name, description string) (*CourseType, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrCourseTypeNameRequired
	}

	now := time.Now()
	return &CourseType{
		id:          uuid.New().String(),
		name:        name,
		description: description,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// Reconstitute 从仓储层还原聚合根。
func Reconstitute(
	id string,
	name string,
	description string,
	createdAt, updatedAt time.Time,
) *CourseType {
	return &CourseType{
		id:          id,
		name:        name,
		description: description,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// --- 核心领域行为 (Domain Behaviors) ---

// UpdateInfo 修改名称与描述。
func (ct *CourseType) UpdateInfo(name, description string) error {
	if strings.TrimSpace(name) == "" {
		return ErrCourseTypeNameRequired
	}
	ct.name = name
	ct.description = description
	ct.updatedAt = time.Now()
	return nil
}

// --- 只读属性访问器 (Getters) ---

func (ct *CourseType) ID() string           { return ct.id }
func (ct *CourseType) Name() string         { return ct.name }
func (ct *CourseType) Description() string  { return ct.description }
func (ct *CourseType) CreatedAt() time.Time { return ct.createdAt }
func (ct *CourseType) UpdatedAt() time.Time { return ct.updatedAt }
