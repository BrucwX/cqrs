package model

import (
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
)

// CourseType 是 Redis 快照（对应 course_type），领域模型是 courseType.CourseType。
//
// 课程类型描述「教什么」（少儿编程 / 成人英语），讲师的授课资质绑在这一层，
// 不绑到具体课程。三个字段直接对应聚合，没有值对象需要打平。
type CourseType struct {
	ID          string    `json:"id"`          // 课程类型 ID（uuid）
	Name        string    `json:"name"`        // 类型名称
	Description string    `json:"description"` // 类型描述
	CreatedAt   time.Time `json:"created_at"`  // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`  // 更新时间
}

// CourseTypeToRedis 写路径。
func CourseTypeToRedis(do *courseType.CourseType) (*CourseType, error) {
	if do == nil {
		return nil, ErrCourseTypeToRedis
	}
	return &CourseType{
		ID:          do.ID(),
		Name:        do.Name(),
		Description: do.Description(),
		CreatedAt:   do.CreatedAt(),
		UpdatedAt:   do.UpdatedAt(),
	}, nil
}

// CourseTypeFromRedis 读路径。
func CourseTypeFromRedis(po *CourseType) (*courseType.CourseType, error) {
	if po == nil {
		return nil, ErrCourseTypeFromRedis
	}
	return courseType.Reconstitute(
		po.ID, po.Name, po.Description, po.CreatedAt, po.UpdatedAt,
	), nil
}
