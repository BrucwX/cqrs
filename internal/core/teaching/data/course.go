package data

import (
	"context"

	"github.com/go-kratos/kratos-layout/internal/core/teaching/biz/course"
	"github.com/go-kratos/kratos-layout/internal/shared/types"

	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type courseRepo struct {
	data *Data
}

// NewCourseRepo creates a new CourseRepo instance.
func NewCourseRepo(d *Data) course.CourseRepo {
	return &courseRepo{data: d}
}

func (r *courseRepo) FindByID(ctx context.Context, id string) (*course.Course, error) {
	// TODO: implement database query
	return nil, nil
}

func (r *courseRepo) ListCourses(ctx context.Context, opts ...types.ListOption) ([]*course.Course, error) {
	// TODO: implement database query
	return nil, nil
}

func (r *courseRepo) CreateCourse(ctx context.Context, c *course.Course) (*course.Course, error) {
	// TODO: implement database insert
	return c, nil
}

func (r *courseRepo) UpdateCourse(ctx context.Context, c *course.Course, mask *fieldmaskpb.FieldMask) (*course.Course, error) {
	// TODO: implement database update
	// 根据 mask 决定更新哪些字段
	return c, nil
}

func (r *courseRepo) DeleteCourse(ctx context.Context, id string) error {
	// TODO: implement database soft delete
	return nil
}
