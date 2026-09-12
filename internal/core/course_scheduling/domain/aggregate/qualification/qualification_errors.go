package qualification

import "errors"

var (
	// ErrTeacherRequired 授予资质时必须给出讲师 ID。
	ErrTeacherRequired = errors.New("teacher is required")
	// ErrCourseTypeRequired 授予资质时必须给出课程类型 ID。
	ErrCourseTypeRequired = errors.New("course type is required")

	ErrQualificationRevoked = errors.New("teacher qualification for this course has been revoked")

	// ErrQualificationRequired 传入的授课资质为空。
	ErrQualificationRequired = errors.New("qualification is required")
	// ErrQualificationNotFound 指定的授课资质不存在。
	ErrQualificationNotFound = errors.New("qualification not found")
	// ErrCourseNotFinished 讲师还没以学员身份修完该课程类型下的课程（没结业，或缺过课）。
	ErrCourseNotFinished = errors.New("teacher has not finished a course of this type")
)
