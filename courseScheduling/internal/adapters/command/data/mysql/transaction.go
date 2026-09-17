package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-kratos/kratos/v3/log"
)

// txKey 是 context 里携带事务的 key。
type txKey struct{}

// txValue 是 context 里携带的事务：既能当连接用（Querier），也能提交/回滚。
type txValue struct {
	Querier
	tx *sql.Tx
}

// Conn 返回本次调用该用的连接：ctx 里带着事务就用事务，否则用连接池。
func (d *MysqlData) Conn(ctx context.Context) Querier {
	if tx, ok := ctx.Value(txKey{}).(txValue); ok {
		return tx
	}
	return Debug(d.db, d.debug)
}

// Begin 开启事务，返回带着它的 ctx。
//
// 之后必须用这个 ctx 去调仓库，否则 d.Conn(ctx) 拿到的还是连接池 ——
// 「先读排期判冲突、再写回」能原子完成的前提，就是判定读到的与写回用的是同一个
// 事务，否则判定读到的会是事务外的快照。
func (d *MysqlData) Begin(ctx context.Context) (context.Context, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("mysql: begin: %w", err)
	}
	return context.WithValue(ctx, txKey{}, txValue{Querier: Debug(tx, d.debug), tx: tx}), nil
}

// End 结束事务：err 非 nil 就回滚，否则提交。
//
// 必须传 Begin 返回的那个 ctx。传错 ctx 时报错而不是默默什么都不做 ——
// 那意味着事务根本没开，写进去的东西也不会提交。
func (d *MysqlData) End(ctx context.Context, err error) error {
	txv, ok := ctx.Value(txKey{}).(txValue)
	if !ok {
		return fmt.Errorf("mysql: end: no transaction in context")
	}

	if err != nil {
		if rbErr := txv.tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			log.Error("mysql rollback", "err", rbErr)
		}
		return err
	}
	return txv.tx.Commit()
}
