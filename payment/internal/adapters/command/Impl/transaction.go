package command

import (
	"context"

	appcmd "payment/internal/app/command/appRepo"
)

// Begin 开启事务，转发给 MysqlData。
func (d *CommandImpl) Begin(ctx context.Context) (context.Context, error) {
	return d.MysqlData.Begin(ctx)
}

// Commit 提交事务，转发给 MysqlData。
func (d *CommandImpl) Commit(ctx context.Context) error {
	return d.MysqlData.Commit(ctx)
}

// Rollback 回滚事务，转发给 MysqlData。
func (d *CommandImpl) Rollback(ctx context.Context) error {
	return d.MysqlData.Rollback(ctx)
}

// 编译期约束：CommandImpl 必须实现 app/command.Transaction。
var _ appcmd.Transaction = (*CommandImpl)(nil)
