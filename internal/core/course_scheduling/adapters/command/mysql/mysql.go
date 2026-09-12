// Package mysql 是 course_scheduling 写侧的 MySQL 实现。
//
// 写侧连的是「写库」（主库）：
//
//	data.database.source  写库连接串
//
// 读侧（adapters/query/mysql）连的是 data.database.read_source，两者可以不是
// 同一个库，所以两侧各带一套自己的连接与 model，互不共用。
//
//	mysql.go    本文件：写侧连接（Data / NewData / Conn / WithTx）
//	driver.go   连库底座：Querier 抽象、DSN 归一化、Open、SQL 日志
//
// model/ 与 implement/ 等有真实实现时再补：届时实现 domain/repo/command
// 的 16 个接口，对应读侧的 model/ 与 implement/。
//
// 与读侧的关键差别是事务：写侧提供 WithTx，把事务放进 context，
// 这样「先读排期判冲突、再写回」里那些经仓库发起的读也能落在同一个事务上。
// ⚠️ 前提是读也走写库；如果读走 adapters/query/mysql 的从库，就不构成原子读了。
package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cqrs/internal/conf"

	"github.com/go-kratos/kratos/v3/log"
)

// Data 持有写库连接池，写侧的所有仓库共享它。
type Data struct {
	db    *sql.DB
	debug bool
}

// NewData 打开写库连接池。
func NewData(c *conf.Data) (*Data, func(), error) {
	cfg := c.GetDatabase()
	if cfg == nil {
		return nil, nil, errors.New("mysql/command: data.database is required")
	}
	if err := ValidateDriver(cfg); err != nil {
		return nil, nil, err
	}

	db, cleanup, err := Open(cfg.GetSource())
	if err != nil {
		return nil, nil, err
	}
	return &Data{db: db, debug: cfg.GetDebug()}, cleanup, nil
}

// Conn 返回本次调用该用的连接：ctx 里带着事务就用事务，否则用连接池。
func (d *Data) Conn(ctx context.Context) Querier {
	if tx, ok := ctx.Value(txKey{}).(Querier); ok {
		return tx
	}
	return Debug(d.db, d.debug)
}

// txKey 是 context 里携带事务的 key。
type txKey struct{}

// WithTx 在一个事务里执行 fn，出错自动回滚。
//
// fn 收到的 ctx 携带该事务，因此期间任何 d.Conn(ctx) 拿到的连接都是同一个事务 ——
// 这正是「先读排期判冲突、再写回」能以原子方式完成的前提，
// 否则判定读到的会是事务外的快照。
func (d *Data) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("mysql: begin: %w", err)
	}

	q := Debug(tx, d.debug)
	if err := fn(context.WithValue(ctx, txKey{}, q)); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			log.Error("mysql rollback", "err", rbErr)
		}
		return err
	}
	return tx.Commit()
}
