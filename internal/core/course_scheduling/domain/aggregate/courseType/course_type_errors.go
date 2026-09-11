package courseType

import "errors"

var (
	// ErrCourseTypeNameRequired 课程类型名称不能为空。
	ErrCourseTypeNameRequired = errors.New("course type name is required")
)
