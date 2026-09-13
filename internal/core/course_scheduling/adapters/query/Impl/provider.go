package query

import (
	"github.com/google/wire"

	"cqrs/internal/core/course_scheduling/adapters/query/data/mysql"
	repoquery "cqrs/internal/core/course_scheduling/domain/repo/query"
)

// ProviderSet 是读侧查询适配器的 provider 集合：存储客户端 + 载体 + 接口绑定。
//
// QueryImpl 一个载体实现全部 10 个 query 接口，这里把它们一次绑好。
var ProviderSet = wire.NewSet(
	mysql.NewData,
	NewQueryImpl,
	wire.Bind(new(repoquery.AbsenceRecordQuery), new(*QueryImpl)),
	wire.Bind(new(repoquery.ClassroomQuery), new(*QueryImpl)),
	wire.Bind(new(repoquery.CourseQuery), new(*QueryImpl)),
	wire.Bind(new(repoquery.CourseSlotQuery), new(*QueryImpl)),
	wire.Bind(new(repoquery.CourseSlotChangeQuery), new(*QueryImpl)),
	wire.Bind(new(repoquery.CourseEnrollmentQuery), new(*QueryImpl)),
	wire.Bind(new(repoquery.StudentMakeupQuery), new(*QueryImpl)),
	wire.Bind(new(repoquery.QualificationQuery), new(*QueryImpl)),
	wire.Bind(new(repoquery.StudentQuery), new(*QueryImpl)),
	wire.Bind(new(repoquery.TeacherQuery), new(*QueryImpl)),
)
