package courseType

import "errors"

var (
	// ErrCourseTypeNameRequired 课程类型名称不能为空。
	ErrCourseTypeNameRequired = errors.New("course type name is required")
	// ErrCourseTypeNotFound 课程或它归属的课程类型不存在。
	ErrCourseTypeNotFound = errors.New("course type not found")
)
