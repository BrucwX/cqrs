package absence

import "errors"

var (
	ErrInvalidHours = errors.New("missed hours must be greater than zero")

	// ErrAbsenceRequired 传入的缺勤记录为空。
	ErrAbsenceRequired = errors.New("absence record is required")
	// ErrAbsenceNotFound 指定的缺勤记录不存在。
	ErrAbsenceNotFound = errors.New("absence record not found")
)
