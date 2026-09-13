package command

import (
	"context"

	appcmd "cqrs/internal/core/course_scheduling/app/command"
)

// Begin 开启事务，转发给 MysqlData。
func (d *CommandImpl) Begin(ctx context.Context) (context.Context, error) {
	return d.MysqlData.Begin(ctx)
}

// End 结束事务，转发给 MysqlData。
func (d *CommandImpl) End(ctx context.Context, err error) error {
	return d.MysqlData.End(ctx, err)
}

// 编译期约束：CommandImpl 必须实现 app/command.Transaction。
var _ appcmd.Transaction = (*CommandImpl)(nil)
