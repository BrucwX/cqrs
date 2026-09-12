package classroomCapacity

import "cqrs/internal/core/course_scheduling/domain/repo/command"

// Service 教室容量相容判定。
//
// 两种问法，问的都是同一件事 ——「这间教室装得下这节课要点名的人吗」，
// 区别只在人数怎么算、教室从哪来：
//
//	Check        给一批槽位挑教室：教室由调用方给，逐槽比「该课人数上限」
//	CheckMakeup  学员要补到某节课上：教室从目标槽位里取，人数还要加上
//	             同节课上其他补课学员（课程上限 + 其他补课学员 + 这位）
//
// 跨了「教室 × 排期 × 课程（× 补课预约）」几个聚合根：槽位只说它属于哪门课，
// 一门课收多少人记在课程上，教室能装多少记在教室上 —— 单个聚合看不全，
// 所以独立成服务。
//
// 只管容量。教室是不是在用、是不是已经被占了座位由教室聚合自己的
// CanAccommodate 一起判；时间撞不撞是 schedule.Conflict 的事。
//
// 返回 true 表示「装不下」，调用方应当拒绝本次安排 —— 与 schedule.Conflict /
// qualification.Check 一致：true = 拦下。
type Service struct {
	classrooms command.ClassroomCommand
	slots      command.CourseSlotCommand
	courses    command.CourseCommand
	makeups    command.StudentMakeupCommand
}

// NewService 创建教室容量相容判定服务。
func NewService(
	classrooms command.ClassroomCommand,
	slots command.CourseSlotCommand,
	courses command.CourseCommand,
	makeups command.StudentMakeupCommand,
) *Service {
	return &Service{classrooms: classrooms, slots: slots, courses: courses, makeups: makeups}
}
