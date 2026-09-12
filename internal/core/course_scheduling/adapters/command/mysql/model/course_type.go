package model

import (
	"time"

	"cqrs/internal/core/course_scheduling/domain/aggregate/courseType"
)

// CourseType 对应表 course_type，领域模型是 courseType.CourseType。
//
// 课程类型描述「教什么」（少儿编程 / 成人英语），讲师的授课资质绑在这一层，
// 不绑到具体课程。三个字段直接对应聚合，没有值对象需要打平。
type CourseType struct {
	ID          string    // varchar(36) 课程类型 ID（uuid）
	Name        string    // varchar(64)  类型名称
	Description string    // varchar(255) 类型描述
	CreatedAt   time.Time // datetime     创建时间
	UpdatedAt   time.Time // datetime     更新时间
	LockVersion uint64    // bigint unsigned 乐观锁版本号
}

// CourseTypeToPO 写路径。
func CourseTypeToPO(do *courseType.CourseType) *CourseType {
	return &CourseType{
		ID:          do.ID(),
		Name:        do.Name(),
		Description: do.Description(),
		CreatedAt:   do.CreatedAt(),
		UpdatedAt:   do.UpdatedAt(),
	}
}

// CourseTypeToDO 读路径。
func CourseTypeToDO(po *CourseType) (*courseType.CourseType, error) {
	return courseType.Reconstitute(
		po.ID, po.Name, po.Description, po.CreatedAt, po.UpdatedAt,
	), nil
}
