package command

import (
	"context"

	"payment/internal/ports"
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

// 编译期约束：CommandImpl 必须实现 ports.Transaction。
var _ ports.Transaction = (*CommandImpl)(nil)
