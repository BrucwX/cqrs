// Package appRepo 定义命令侧的存储端口。
//
// 单独放在一个子包里，而不是和 provider.go 一起放 app/command 父包：
// handler 子包要 import 这个端口，而父包的 provider.go 又要 import handler 子包，
// 两端都不碰对方，才不会成环。
package appRepo

import "context"

// Transaction 手动开/关一个事务。
//
type Transaction interface {
	// Begin 开启事务，返回带着它的 ctx。
	Begin(ctx context.Context) (context.Context, error)

	// Commit 提交事务（传 Begin 返回的 ctx）。
	Commit(ctx context.Context) error

	// Rollback 回滚事务（传 Begin 返回的 ctx）。
	Rollback(ctx context.Context) error
}
