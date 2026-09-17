// Package mysql 是 course_scheduling 读侧的 MySQL 实现。
//
// 读侧只干一件事：把数据查出来。它连的是「读库」——
//
//	data.database.read_source  读库连接串（留空则回落到 source，单库环境只配一个就够）
//	data.database.source       写库连接串，命令侧（adapters/command/data/mysql）用
//
// 读库和写库可以是同一个库，也可以是主从、甚至两套不同的库，所以读侧自带
// 自己的一份 model 与连接，不和写侧共用任何东西。
//
//	mysql.go    本文件：读库连接池（Data / NewData / Conn）
//	driver.go   连库底座：Querier 抽象、DSN 归一化、Open、SQL 日志
//	<资源>.go   各聚合根的查询实现（Data 直接实现全部 domain/repo/query 接口）
//	model/      数据模型（PO）+ DO <-> PO 转换函数
//	help/       Scanner 抽象 + 列清单 + 分页/扫行辅助
//
// ⚠️ 读侧刻意不提供事务。如果 read_source 指向从库，那么命令侧在 Begin/End
// 之间通过查询仓库发起的读会落到从库的快照上，与写事务不构成原子读 ——
// 需要「读判定 + 写回」原子的场景，读要走在写库上（把 read_source 留空即可）。
package mysql

import (
	"context"
	"database/sql"
	"errors"

	"cqrs/internal/conf"
)

// Data 持有读库连接池，读侧的所有仓库共享它。
type Data struct {
	db    *sql.DB
	debug bool
}

// NewData 打开读库连接池。
func NewData(c *conf.Data) (*Data, func(), error) {
	cfg := c.GetDatabase()
	if cfg == nil {
		return nil, nil, errors.New("mysql/query: data.database is required")
	}
	if err := ValidateDriver(cfg); err != nil {
		return nil, nil, err
	}

	db, cleanup, err := Open(readSource(cfg))
	if err != nil {
		return nil, nil, err
	}
	return &Data{db: db, debug: cfg.GetDebug()}, cleanup, nil
}

// readSource 挑出读库连接串：优先 read_source，没配就用 source。
func readSource(cfg *conf.Data_Database) string {
	if source := cfg.GetReadSource(); source != "" {
		return source
	}
	return cfg.GetSource()
}

// Conn 返回查询该用的连接。
//
// 读侧不开事务，所以永远返回连接池本身。
func (d *Data) Conn(_ context.Context) Querier {
	return Debug(d.db, d.debug)
}
