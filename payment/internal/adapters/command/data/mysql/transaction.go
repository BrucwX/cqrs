package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

// Commit 提交事务。必须传 Begin 返回的那个 ctx。
func (d *MysqlData) Commit(ctx context.Context) error {
	txv, ok := ctx.Value(txKey{}).(txValue)
	if !ok {
		return fmt.Errorf("mysql: commit: no transaction in context")
	}
	return txv.tx.Commit()
}

// Rollback 回滚事务。必须传 Begin 返回的那个 ctx。
//
// 已经结束的事务（sql.ErrTxDone）不算错 —— defer 里重复回滚不该再报一次。
func (d *MysqlData) Rollback(ctx context.Context) error {
	txv, ok := ctx.Value(txKey{}).(txValue)
	if !ok {
		return fmt.Errorf("mysql: rollback: no transaction in context")
	}
	if err := txv.tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
		return err
	}
	return nil
}
