package command

import (
	"github.com/google/wire"

	"cqrs/internal/core/course_scheduling/adapters/command/data/mysql"
	appcmd "cqrs/internal/core/course_scheduling/app/command/appRepo"
	repo "cqrs/internal/core/course_scheduling/domain/repo/command"
)

// ProviderSet 是写侧命令适配器的 provider 集合：存储客户端 + 载体 + 接口绑定。
//
// CommandImpl 一个载体实现全部 11 个 command 接口 + Transaction，
// 这里把它们一次绑好；调用方（course_scheduling.ProviderSet）只引用这一个集合。
var ProviderSet = wire.NewSet(
	mysql.NewMysqlData,
	NewCommandImpl,
	wire.Bind(new(appcmd.Transaction), new(*CommandImpl)),
	wire.Bind(new(repo.AbsenceRecordCommand), new(*CommandImpl)),
	wire.Bind(new(repo.ClassroomCommand), new(*CommandImpl)),
	wire.Bind(new(repo.CourseCommand), new(*CommandImpl)),
	wire.Bind(new(repo.CourseSlotCommand), new(*CommandImpl)),
	wire.Bind(new(repo.CourseSlotChangeCommand), new(*CommandImpl)),
	wire.Bind(new(repo.CourseTypeCommand), new(*CommandImpl)),
	wire.Bind(new(repo.CourseEnrollmentCommand), new(*CommandImpl)),
	wire.Bind(new(repo.StudentMakeupCommand), new(*CommandImpl)),
	wire.Bind(new(repo.QualificationCommand), new(*CommandImpl)),
	wire.Bind(new(repo.StudentCommand), new(*CommandImpl)),
	wire.Bind(new(repo.TeacherCommand), new(*CommandImpl)),
)
