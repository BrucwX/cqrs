package help

import (
	"context"
	"database/sql"
)

// Querier 是 help 需要的最小查询接口：只读，不 import mysql，避免环。
type Querier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// List 执行查询，把每一行都交给 scanRow 转成 PO。
//
// 返回空切片而不是 nil，方便调用方直接序列化。
func List[T any](
	ctx context.Context,
	q Querier,
	query string,
	args []any,
	scanRow func(Scanner) (*T, error),
) ([]*T, error) {
	rows, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	out := make([]*T, 0)
	for rows.Next() {
		item, err := scanRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// convertAll 把一批 PO 逐条转成 DO。
//
// 读路径的转换会做值对象校验，任何一条脏数据都会让整次查询失败 ——
// 宁可报错也不要返回一个不满足不变量的聚合。
func convertAll[PO, DO any](pos []*PO, toDO func(*PO) (*DO, error)) ([]*DO, error) {
	out := make([]*DO, 0, len(pos))
	for _, po := range pos {
		item, err := toDO(po)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}

// QueryAll = List + convertAll，查询侧最常见的一步到位写法。
func QueryAll[PO, DO any](
	ctx context.Context,
	q Querier,
	query string,
	args []any,
	scanRow func(Scanner) (*PO, error),
	toDO func(*PO) (*DO, error),
) ([]*DO, error) {
	pos, err := List(ctx, q, query, args, scanRow)
	if err != nil {
		return nil, err
	}
	return convertAll(pos, toDO)
}

// Exists 执行一条 SELECT EXISTS(...)，返回是否命中。
func Exists(ctx context.Context, q Querier, query string, args ...any) (bool, error) {
	var found bool
	if err := q.QueryRowContext(ctx, query, args...).Scan(&found); err != nil {
		return false, err
	}
	return found, nil
}
