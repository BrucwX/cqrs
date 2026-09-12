package command

import (
	"context"
	"errors"

	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
)

var (
	// ErrQualificationRequired 传入的授课资质为空。
	ErrQualificationRequired = errors.New("qualification is required")
	// ErrQualificationNotFound 指定的授课资质不存在。
	ErrQualificationNotFound = errors.New("qualification not found")
	// ErrCourseNotFinished 讲师还没以学员身份修完该课程类型下的课程（没结业，或缺过课）。
	ErrCourseNotFinished = errors.New("teacher has not finished a course of this type")
)

// QualificationCommand 授课资质命令接口
//
// 资质只「发」和「撤」：能不能拿到由领域规则决定（见 domain/service/qualification
// 的 Check.CheckTeacherGrantable），不存在任人填字段的通用更新。
type QualificationCommand interface {
	// GrantQualification 授予授课资质
	//
	// 只管写：把这条资质落库。够不够格（讲师修没修完该类型的课）由调用方在
	// 调过来之前判完（见 domain/service/qualification 的 Check.CheckTeacherGrantable），所以这里没有
	// 回调，也没有「传 nil 表示不检查」这类分支。
	GrantQualification(ctx context.Context, q *qualification.Qualification) error
	// Delete 删除授课资质
	Delete(ctx context.Context, id int64) error
	// GetQualifications 取该讲师持有的全部资质
	GetQualifications(ctx context.Context, teacherID int64) ([]qualification.Qualification, error)
}
