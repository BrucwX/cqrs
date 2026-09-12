package scheduleConflict

import (
	"cqrs/internal/core/course_scheduling/domain/repo/command"
)

// Service 排期冲突判定。
//
// 六种场景问的其实是同一个问题 —— 「要占的这段时间，是不是已经有人占了」，
// 区别只在「要占的」和「已有的」各自从哪来：
//
//	CheckTeacher     要排给某讲师的那批槽位 vs 该讲师现有的课
//	CheckClassroom   要排给某教室的那批槽位 vs 该教室现有的课
//	CheckCourse      要排给某课程的那批槽位 vs 该课程现有的课
//	CheckStudent     学员要选的那门课     vs 该学员在学课程的排期
//	CheckEnrollment  选课：先过课程自己的准入（窗口 / 容量），再问时间
//	CheckSlotChange  临时换课的目标时段   vs 讲师 / 教室的课表 + 其他换课单
//
// 前四者是「资源现有的排期 vs 本次要排的槽位」；CheckStudent 比的是
// 「学员在学课程的排期 vs 目标课程的排期」—— 学员名下没有槽位，占用集是从他的
// 报名记录推出来的，所以第二个入参给的是课程 ID 而不是槽位 ID。
//
// CheckEnrollment 是 CheckStudent 的完整版：选课除了不能撞时间，还得课程还开着
// 选课窗口、还没满员 —— 那两条是课程聚合自己的规则，一并在这里过完，调用方拿到
// 的就是「能不能选」这一个答案。
//
// CheckSlotChange 比的是「还没落成排期的换课单」：占用既有已经排好的课，也有刚
// 登记还没生效的换课，只看一张表看不全。
//
// 只判这些。讲师有没有资质、教室装不装得下这类事归各自聚合（或别的服务）管，
// 不在这里（所以判定完，调用方还得自己过一遍那些规则）。
//
// 数据来自两个聚合根：已经排好的课都在 CourseSlot 上（CourseSlotCommand），
// 还没落成排期的换课单在 CourseSlotChange 上（CourseSlotChangeCommand）。
// teachers / courses 只用来校入参里的资源确实存在 —— 传错 ID 应该报出来，
// 而不是当成「没冲突」悄悄放行。
//
// 返回 true 表示「有冲突」/「不能选」，调用方应当拒绝这次操作。
type Service struct {
	SlotCmd  command.CourseSlotCommand
	teachers command.TeacherCommand
	courses  command.CourseCommand
	changes  command.CourseSlotChangeCommand
}

// NewService 创建排期冲突判定服务。
func NewService(
	slots command.CourseSlotCommand,
	teachers command.TeacherCommand,
	courses command.CourseCommand,
	changes command.CourseSlotChangeCommand,
) *Service {
	return &Service{SlotCmd: slots, teachers: teachers, courses: courses, changes: changes}
}
