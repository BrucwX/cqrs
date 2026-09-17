package model

import "errors"

// 转换失败哨兵错误：每个聚合一组，DO -> Redis 与 Redis -> DO 各一个。
//
// XxxToRedis 在入参 DO 为 nil 时返回 ToRedis；XxxFromRedis 在入参模型为 nil 或
// 值对象校验不过时返回 FromRedis（后者用 %w 包住领域模型的底层错误）。
var (
	ErrAbsenceToRedis            = errors.New("absence: DO -> Redis failed")
	ErrAbsenceFromRedis          = errors.New("absence: Redis -> DO failed")
	ErrClassroomToRedis          = errors.New("classroom: DO -> Redis failed")
	ErrClassroomFromRedis        = errors.New("classroom: Redis -> DO failed")
	ErrCourseToRedis             = errors.New("course: DO -> Redis failed")
	ErrCourseFromRedis           = errors.New("course: Redis -> DO failed")
	ErrCourseSlotToRedis         = errors.New("course slot: DO -> Redis failed")
	ErrCourseSlotFromRedis       = errors.New("course slot: Redis -> DO failed")
	ErrCourseSlotChangeToRedis   = errors.New("course slot change: DO -> Redis failed")
	ErrCourseSlotChangeFromRedis = errors.New("course slot change: Redis -> DO failed")
	ErrCourseTypeToRedis         = errors.New("course type: DO -> Redis failed")
	ErrCourseTypeFromRedis       = errors.New("course type: Redis -> DO failed")
	ErrEnrollmentToRedis         = errors.New("enrollment: DO -> Redis failed")
	ErrEnrollmentFromRedis       = errors.New("enrollment: Redis -> DO failed")
	ErrMakeupToRedis             = errors.New("makeup: DO -> Redis failed")
	ErrMakeupFromRedis           = errors.New("makeup: Redis -> DO failed")
	ErrQualificationToRedis      = errors.New("qualification: DO -> Redis failed")
	ErrQualificationFromRedis    = errors.New("qualification: Redis -> DO failed")
	ErrStudentToRedis            = errors.New("student: DO -> Redis failed")
	ErrStudentFromRedis          = errors.New("student: Redis -> DO failed")
	ErrTeacherToRedis            = errors.New("teacher: DO -> Redis failed")
	ErrTeacherFromRedis          = errors.New("teacher: Redis -> DO failed")
)
