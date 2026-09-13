// Package mysql 是 course_scheduling 写侧的 MySQL 实现。
//
// 写侧连的是「写库」（主库）：
//
//	data.database.source  写库连接串
//
// 读侧（adapters/query/data/mysql）连的是 data.database.read_source，两者可以不是
// 同一个库，所以两侧各带一套自己的连接与 model，互不共用。
//
//	mysql.go        本文件：写库连接池（MysqlData / NewMysqlData）
//	transaction.go  事务与连接选择（Conn / Begin / End）
//	driver.go       连库底座：Querier 抽象、DSN 归一化、Open、SQL 日志
//	<资源>.go       各聚合根的 repo 实现（MysqlData 直接实现全部接口）
//	model/          PO 与 DO 转换（与读侧逐字节相同）
//	help/           Scanner 抽象 + 列清单 + 扫行辅助
//
// 与读侧的关键差别是事务：写侧提供 Begin / End，把事务放进 context，
// 这样「先读排期判冲突、再写回」里那些经仓库发起的读也能落在同一个事务上。
// ⚠️ 前提是读也走写库；如果读走 adapters/query/data/mysql 的从库，就不构成原子读了。
package mysql

import (
	"database/sql"
	"errors"

	"cqrs/internal/conf"
)

// MysqlData 持有写库连接池，写侧的所有仓库共享它。
type MysqlData struct {
	db    *sql.DB
	debug bool
}

// NewMysqlData 打开写库连接池。
func NewMysqlData(c *conf.Data) (*MysqlData, func(), error) {
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
	return &MysqlData{db: db, debug: cfg.GetDebug()}, cleanup, nil
}
