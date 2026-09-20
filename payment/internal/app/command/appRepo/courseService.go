package appRepo

import "context"

// CourseService 是「课程那一侧」的能力端口：由适配器实现，转去调 courseScheduling。
type CourseService interface {
	// CheckAndEnrollCourse 校验并报名课程：accountIDs 是这批人的**登录账号 id**，
	// courseIDs 是这批商品指向的课程（Products.TargetIds()）。
	//
	// 两边都是 string：跨上下文传的 id 统一用字符串，不在这层暴露各自的数字主键。
	CheckAndEnrollCourse(ctx context.Context, accountIDs []string, courseIDs []string) (err error)
}
