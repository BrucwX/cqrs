package qualificationCheck

import "cqrs/internal/core/course_scheduling/domain/repo/command"

// Service 授课资质判定。
//
// 两个方向的问题都归它，问的都是「这位讲师的资质」：
//
//	CheckTeacherQualifiedForCourseType / CheckTeacherQualified / CheckTeacherQualifiedForSlot
//	    他手上有没有资质 —— 能不能排课。资质本身就绑在「讲师 + 课程类型」上，
//	    三个入口只是入参不一样（课程类型 / 课程 / 槽位），一层层往里解析。
//	CheckTeacherGrantable
//	    该不该给他发这张证 —— 得先以学员身份修完该类型下的一门课。
//
// 一次只判一个对象。要问「这批槽位是不是门门都有资质」由调用方逐个问 ——
// 这样「有一个不过就整批拦下」的决定权留在调用方手里。
//
// 只管资质本身，不看时间：解析时只取「属于哪门课 / 哪个类型」，不碰 weekday /
// 时间区间，那些是 schedule.Conflict 的事。
//
// 为什么要独立成服务：资质挂在讲师身上，课程类型是课程的属性，「这位讲师能不能教
// 这门课」跨了讲师与课程两个聚合根，单个聚合看不全；发证那条还跨了报名与缺勤。
//
// 返回 true 表示「没资质」/「不够格」，调用方应当拒绝本次操作 ——
// 与 schedule.Conflict 一致：true = 拦下。
type Service struct {
	slots          command.CourseSlotCommand
	teachers       command.TeacherCommand
	qualifications command.QualificationCommand
	courseTypes    command.CourseTypeCommand

	// 下面三个只有发证那条用得到
	courses     command.CourseCommand
	enrollments command.CourseEnrollmentCommand
	absences    command.AbsenceRecordCommand
}

// NewService 创建授课资质判定服务。
func NewService(
	slots command.CourseSlotCommand,
	teachers command.TeacherCommand,
	qualifications command.QualificationCommand,
	courseTypes command.CourseTypeCommand,
	courses command.CourseCommand,
	enrollments command.CourseEnrollmentCommand,
	absences command.AbsenceRecordCommand,
) *Service {
	return &Service{
		slots:          slots,
		teachers:       teachers,
		qualifications: qualifications,
		courseTypes:    courseTypes,
		courses:        courses,
		enrollments:    enrollments,
		absences:       absences,
	}
}
