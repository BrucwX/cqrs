// Package model 是 course_scheduling 写侧的数据模型。
//
// 定位：把 script/mysql/course_scheduling.sql 里的表，一对一映射成 Go 结构体（PO），
// 一张表一个文件一个结构体，字段名与列名一一对应，列清单与顺序与建表脚本一致。
//
// 本包只有两样东西，而且都在同一份表文件里：
//
//	<表名>.go     PO 结构体 + 它自己的 XxxToPO / XxxToDO 两个转换函数
//	helpers.go    跨表复用的转换零件（可空时间列、DayTime 列）
//
// 不碰 SQL：语句在 implement 包里。
//
// # 与领域模型的关系
//
// 领域模型在 domain/aggregate：聚合根的字段不可导出，只能通过 getter 读、
// 通过构造函数 / Reconstitute 写，值对象（Location、Capacity…）也是独立的类型。
// 本包反过来：**导出字段 + 值对象就地打平**，这样 database/sql 的 Scan 与赋值
// 可以直接对着字段做，不需要任何中间层。
//
// 两个形状之间的互转就是各表文件末尾的那两个函数（XxxToPO 写路径、XxxToDO
// 读路径），读路径会走一遍领域模型的值对象构造函数，非法数据在那里被拦下。
//
// # 表 <-> 结构体 <-> 领域聚合
//
//	course_type         CourseType         courseType.CourseType
//	classroom           Classroom          classroom.Classroom
//	student             Student            student.Student
//	teacher             Teacher            teacher.Teacher
//	course              Course             course.Course
//	course_slot         CourseSlot         courseSlot.CourseSlot
//	course_slot_change  CourseSlotChange   courseSlotChange.CourseSlotChange
//	course_enrollment   CourseEnrollment   enrollment.CourseEnrollment
//	absence_record      AbsenceRecord      absence.AbsenceRecord
//	student_makeup      StudentMakeup      makeup.StudentMakeup
//	qualification       Qualification      qualification.Qualification
//
// # 类型对照
//
//	varchar(n)         -> string     （字段注释里保留原始长度）
//	int                -> int
//	bigint             -> int64
//	bigint unsigned    -> uint64     （只有 lock_version）
//	tinyint unsigned   -> uint8      （枚举，取值范围见字段注释）
//	datetime           -> time.Time
//	datetime NULL      -> sql.NullTime
//	date               -> time.Time  （只用到年月日）
//	time               -> string     （"15:04:05"，见下面的说明）
//
// 时间列在 DSN 里靠 parseTime=True 解析成 time.Time；TIME 列驱动以文本返回，
// 所以 start_time / end_time 用字符串承载。
//
// # 值对象打平
//
//	Location          -> building / floor / room
//	Capacity          -> capacity_max / capacity_enrolled
//	EnrollmentWindow  -> enroll_start_at / enroll_end_at / drop_deadline
//	CoursePeriod      -> period_start_at / period_end_at / total_hours / completed_hours
//	ContactInfo       -> phone / email
//	DayTimeRange      -> start_time / end_time
//	OriginalPlan      -> original_slot_id / original_date / original_teacher_id /
//	                     original_classroom_id / original_start_time / original_end_time
//	TargetPlan        -> target_start_at / target_end_at / target_teacher_id / target_classroom_id
//
// # 两处「模型里就对不上」的地方（如实建模，不在数据模型里悄悄修正）
//
//  1. course_slot.id 是 varchar(36) 的 uuid，但 absence_record.course_slot_id、
//     student_makeup.original_slot_id / target_slot_id、course_slot_change.original_slot_id
//     都是 bigint，指不到具体的课表模板上。
//  2. course_slot 的时间是 time 列（"09:00:00"），course_slot_change 的原计划快照
//     却是 varchar(5)（"09:00"），同一个概念两种表示。
//
// # 不加外键
//
// 表之间不建外键，引用关系只用 *_id 列表达：ID 全由应用生成、写入顺序自由，
// 删除也由应用决定，让 MySQL 少一份约束来源；因此本包的结构体之间也没有任何导航属性。
package model
