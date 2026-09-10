package query

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/classroom"
)

// ClassroomQuery 教室查询接口
type ClassroomQuery interface {
	// Page 分页查询教室列表
	Page(ctx context.Context, page int, pageSize int) ([]*classroom.Classroom, error)
}
