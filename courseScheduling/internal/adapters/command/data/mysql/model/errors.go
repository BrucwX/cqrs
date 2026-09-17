package model

import "errors"

// 转换失败哨兵错误：每个聚合一组，DO -> PO 与 PO -> DO 各一个。
//
// XxxToPO 在入参 DO 为 nil 时返回 DOToPO；XxxToDO 在入参 PO 为 nil 或
// 值对象校验不过时返回 POToDO（后者用 %w 包住领域模型的底层错误）。
var (
	ErrAbsenceDOToPO          = errors.New("absence: DO -> PO conversion failed")
	ErrAbsencePOToDO          = errors.New("absence: PO -> DO conversion failed")
	ErrClassroomDOToPO        = errors.New("classroom: DO -> PO conversion failed")
	ErrClassroomPOToDO        = errors.New("classroom: PO -> DO conversion failed")
	ErrCourseDOToPO           = errors.New("course: DO -> PO conversion failed")
	ErrCoursePOToDO           = errors.New("course: PO -> DO conversion failed")
	ErrCourseSlotDOToPO       = errors.New("course slot: DO -> PO conversion failed")
	ErrCourseSlotPOToDO       = errors.New("course slot: PO -> DO conversion failed")
	ErrCourseSlotChangeDOToPO = errors.New("course slot change: DO -> PO conversion failed")
	ErrCourseSlotChangePOToDO = errors.New("course slot change: PO -> DO conversion failed")
	ErrCourseTypeDOToPO       = errors.New("course type: DO -> PO conversion failed")
	ErrCourseTypePOToDO       = errors.New("course type: PO -> DO conversion failed")
	ErrEnrollmentDOToPO       = errors.New("enrollment: DO -> PO conversion failed")
	ErrEnrollmentPOToDO       = errors.New("enrollment: PO -> DO conversion failed")
	ErrMakeupDOToPO           = errors.New("makeup: DO -> PO conversion failed")
	ErrMakeupPOToDO           = errors.New("makeup: PO -> DO conversion failed")
	ErrQualificationDOToPO    = errors.New("qualification: DO -> PO conversion failed")
	ErrQualificationPOToDO    = errors.New("qualification: PO -> DO conversion failed")
	ErrStudentDOToPO          = errors.New("student: DO -> PO conversion failed")
	ErrStudentPOToDO          = errors.New("student: PO -> DO conversion failed")
	ErrTeacherDOToPO          = errors.New("teacher: DO -> PO conversion failed")
	ErrTeacherPOToDO          = errors.New("teacher: PO -> DO conversion failed")
)
