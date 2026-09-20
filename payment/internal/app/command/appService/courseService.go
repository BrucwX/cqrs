package appservice

import "context"

type CourseService interface {
	// MustGetCourseByID 按 ID 查课程：不存在就报错。
	CheckAndEnrollCourse(ctx context.Context, userID int64, productID int64) (err error)
}
