package qualification

import "errors"

var (
	// ErrTeacherRequired 授予资质时必须给出讲师 ID。
	ErrTeacherRequired = errors.New("teacher is required")
	// ErrCourseTypeRequired 授予资质时必须给出课程类型 ID。
	ErrCourseTypeRequired = errors.New("course type is required")

	ErrQualificationRevoked = errors.New("teacher qualification for this course has been revoked")
)
