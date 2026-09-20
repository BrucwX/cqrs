package appRepo

import "context"

// CourseService 是「课程那一侧」的能力端口：由适配器实现，转去调 courseScheduling。
type CourseService interface {
	// 课程那边 ,检查能否注册,注册, 以及（减库存） 
	CheckAndEnrollCourse(ctx context.Context, accountIDs []string, courseIDs []string) (err error)
	// 课程那边注销
	UnenrollCourse(ctx context.Context, accountIDs []string, courseIDs []string) (err error)
}
