// 本文件是写侧连库的底座：怎么开连接池、怎么补 DSN、怎么把语句打到日志里。
//
// 刻意不抽到两侧共用的包里 —— 读侧和写侧本来就是可以连不同库的独立实现，
// 各自带一份最省心。真要 DRY 的话把本文件原样提到一个中立的包即可。
package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"cqrs/internal/conf"

	// 注册 database/sql 的 "mysql" 驱动。
	"github.com/go-kratos/kratos/v3/log"
	_ "github.com/go-sql-driver/mysql"
)

// driverName 是本适配器支持的 database/sql 驱动名。
const driverName = "mysql"

// ValidateDriver 校验配置里的驱动名是不是本适配器支持的。
//
// 留空视为没指定，按 mysql 处理。
func ValidateDriver(cfg *conf.Data_Database) error {
	if driver := cfg.GetDriver(); driver != "" && driver != driverName {
		return fmt.Errorf("mysql: unsupported driver %q, want %q", driver, driverName)
	}
	return nil
}

// Open 打开一个 MySQL 连接池并做一次连通性检查，同时配好连接池参数。
//
// 返回的清理函数由调用方决定什么时候关（wire 里挂在进程退出上）。
func Open(source string) (*sql.DB, func(), error) {
	if strings.TrimSpace(source) == "" {
		return nil, nil, errors.New("mysql: database source is empty")
	}

	db, err := sql.Open(driverName, normalizeDSN(source))
	if err != nil {
		return nil, nil, fmt.Errorf("mysql: open: %w", err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, nil, fmt.Errorf("mysql: ping: %w", err)
	}

	cleanup := func() {
		if err := db.Close(); err != nil {
			log.Error("closing course_scheduling mysql command adapter", "err", err)
			return
		}
		log.Info("closing course_scheduling mysql command adapter")
	}
	return db, cleanup, nil
}

// normalizeDSN 补齐 parseTime / loc / charset。
//
// 时间列全部按 time.Time 读写（DATE / DATETIME 依赖 parseTime=True），
// 少了这个参数驱动会退化成 []byte，扫到 time.Time 直接报错，所以在入口兜住。
//
// 注意 TIME 列（course_slot.start_time / end_time）不受 parseTime 影响，
// 驱动以文本返回，所以数据模型里这两个字段是 string。
func normalizeDSN(dsn string) string {
	base, rawQuery, _ := strings.Cut(dsn, "?")
	params, err := url.ParseQuery(rawQuery)
	if err != nil {
		params = url.Values{}
	}
	defaults := map[string]string{
		"parseTime": "true",
		"loc":       "Local",
		"charset":   "utf8mb4",
	}
	for key, value := range defaults {
		if params.Get(key) == "" {
			params.Set(key, value)
		}
	}
	return base + "?" + params.Encode()
}

// Querier 抽象 *sql.DB 与 *sql.Tx 的公共能力，让同一段 SQL 既能在连接池上跑，
// 也能在事务里跑。
type Querier interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// Debug 按开关决定要不要给连接套一层 SQL 日志壳。
//
// conf.data.database.debug = true 时，每条语句都会打到 debug 日志里。
func Debug(q Querier, on bool) Querier {
	if !on {
		return q
	}
	return debugQuerier{inner: q}
}

// debugQuerier 把每条语句打到 debug 日志里，方便本地排查。
type debugQuerier struct {
	inner Querier
}

func (d debugQuerier) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	log.Debug("mysql exec", "sql", compactSQL(query), "args", truncateArgs(args))
	return d.inner.ExecContext(ctx, query, args...)
}

func (d debugQuerier) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	log.Debug("mysql query", "sql", compactSQL(query), "args", truncateArgs(args))
	return d.inner.QueryContext(ctx, query, args...)
}

func (d debugQuerier) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	log.Debug("mysql query row", "sql", compactSQL(query), "args", truncateArgs(args))
	return d.inner.QueryRowContext(ctx, query, args...)
}

// compactSQL 把多行 SQL 压成一行，日志里更好读。
func compactSQL(query string) string {
	return strings.Join(strings.Fields(query), " ")
}

// truncateArgs 截断过长的字符串参数，避免把整行数据打进日志。
func truncateArgs(args []any) []any {
	const maxLen = 32

	out := make([]any, 0, len(args))
	for _, arg := range args {
		if s, ok := arg.(string); ok && len(s) > maxLen {
			out = append(out, s[:maxLen]+"...")
			continue
		}
		out = append(out, arg)
	}
	return out
}
