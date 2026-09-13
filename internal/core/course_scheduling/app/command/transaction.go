package command

import "context"

// Transaction 手动开/关一个事务。
//
// Begin 返回的 ctx 带着事务，之后必须拿它去调仓库 —— 事务靠 ctx 传递，
// 用外面的 ctx 就跑到事务外了。
//
// 实现在存储侧（见 adapters/command/data/mysql/implement/app.Transaction）：
// 期间任何仓库发出的读也落在同一个事务里。
// ⚠️ 前提是读也走写库。如果读走 adapters/query/mysql 的从库，两边不共享事务，
// 就不构成原子读了。
type Transaction interface {
	// Begin 开启事务，返回带着它的 ctx。
	Begin(ctx context.Context) (context.Context, error)

	// End 结束事务：err 非 nil 就回滚，否则提交。
	//
	// 把 err 原样带出来，好在 defer 里一行收尾（见类型注释里的用法）。
	End(ctx context.Context, err error) error
}
