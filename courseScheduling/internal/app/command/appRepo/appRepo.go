// Package appRepo 定义命令侧的存储端口。
package appRepo

import (
	"context"
	"errors"

	"courseScheduling/internal/domain/aggregate/classroom"
)

// Transaction 真事务：Begin 开一个数据库事务放进 ctx，End 按 err 提交或回滚。
type Transaction interface {
	Begin(ctx context.Context) (context.Context, error)
	End(ctx context.Context, err error) error
}

var ErrVersionConflict = errors.New("optimistic lock version conflict")

// OptimisticTransaction 乐观锁事务：不加锁，靠「读版本 → 写时校验版本」发现并发，
// 版本对不上就重跑**整个用例**。

type OptimisticTransaction interface {
	// RunUpdatClassroom 跑一次「更新教室」：id 是目标教室，updateFn 表达「怎么改」。
	//
	// 读教室（连同 lock_version）、带版本写回、冲突重试都由实现负责：
	// 版本对不上就重跑 updateFn，最多 10 次，用光了才升级为真事务（行锁）。
	// ⚠️ updateFn 会被执行多次，别在里面干发号 / 发消息这类事。
	RunUpdatClassroom(
		ctx context.Context,
		id string,
		updateFn func(ctx context.Context, cl *classroom.Classroom) (*classroom.Classroom, error),
	) error
}
