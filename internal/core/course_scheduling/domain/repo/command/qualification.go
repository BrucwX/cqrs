package command

import (
	"context"

	"cqrs/internal/core/course_scheduling/domain/aggregate/qualification"
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
	// DeleteQualification 删除授课资质
	DeleteQualification(ctx context.Context, id int64) error
	// MustGetQualification 取授课资质聚合；不存在时报 ErrQualificationNotFound
	MustGetQualification(ctx context.Context, id int64) (qualification.Qualification, error)
	// GetQualifications 取该讲师持有的全部资质
	GetQualifications(ctx context.Context, teacherID int64) ([]qualification.Qualification, error)
}
