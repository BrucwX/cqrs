package query

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/teacher"
)

// TeacherQuery 讲师查询接口
type TeacherQuery interface {
	// Page 分页查询讲师列表
	Page(ctx context.Context, page int, pageSize int) ([]*teacher.Teacher, error)
}
