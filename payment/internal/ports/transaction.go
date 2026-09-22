// Package ports 是 payment 的**端口层**：用例（app 层）要用、但不属于任何聚合的外部能力，
// 都在这儿声明成接口。
//
// 目前四类：
//
//	Transaction        事务边界：谁开、谁提交、谁回滚
//	CourseService      别的限界上下文的能力（转调 courseScheduling）
//	IdempotencyService 幂等：这个请求是不是已经办过了
//	PayService         支付渠道：发起 / 查单 / 关单 / 退款
//
// 它们**不是**聚合的仓储（仓储在 domain/ 里、按聚合分）—— 是「聚合之外、但用例绕不开」的事。
//
// 单独一层放顶层：用例只认接口，实现全在 adapters/（MySQL、Redis、渠道 SDK）。
// 依赖箭头一律由外往内指 —— ports 谁都能 import，它自己不 import 任何人。
package ports

import "context"

// Transaction 手动开/关一个事务。
type Transaction interface {
	// Begin 开启事务，返回带着它的 ctx。
	Begin(ctx context.Context) (context.Context, error)

	// Commit 提交事务（传 Begin 返回的 ctx）。
	Commit(ctx context.Context) error

	// Rollback 回滚事务（传 Begin 返回的 ctx）。
	Rollback(ctx context.Context) error
}
