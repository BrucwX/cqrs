package command

import (
	"context"

	appcmd "courseScheduling/internal/app/command/appRepo"
	"courseScheduling/internal/domain/aggregate/classroom"
)

// 编译期约束：CommandImpl 必须实现 appcmd.OptimisticTransaction。
var _ appcmd.OptimisticTransaction = (*CommandImpl)(nil)

// RunUpdatClassroom 转发给 MysqlData：读教室、取版本、带版本写回、冲突重试都在存储侧。
func (d *CommandImpl) RunUpdatClassroom(
	ctx context.Context,
	id string,
	updateFn func(ctx context.Context, cl *classroom.Classroom) (*classroom.Classroom, error),
) error {
	return d.MysqlData.RunUpdatClassroom(ctx, id, updateFn)
}
