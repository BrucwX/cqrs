package query

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/student"
)

// StudentQuery 学员查询接口
type StudentQuery interface {
	// Page 分页查询学员列表
	Page(ctx context.Context, page int, pageSize int) ([]*student.Student, error)
}
