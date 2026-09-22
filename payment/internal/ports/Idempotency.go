package ports

import "errors"

// IdempotencyService 幂等请求端口：按幂等键判断「这个请求是不是已经处理过了」。
//
// ⚠️ 幂等键由**调用方（前端）生成**，同一次提交重试必须复用同一个键 ——
// 后端自己生成的话每次都是新键，等于没做幂等。
//
// 实现建议：键落库 + 唯一约束（先查一次挡常规重试，并发则靠插入撞唯一键兜底）。
type IdempotencyService interface {
	// CheckIdempotency 检查幂等键：没用过返回 nil；用过了返回 ErrDuplicateRequest，
	// 调用方据此回放上次结果（或告诉用户别重复提交）。
	CheckIdempotency(idempotencyKey string) error
}

// ErrDuplicateRequest 这个幂等键已经用过了 —— 重复提交 / 重试，别再办一次事。
//
// 生产上更完整的做法是「回放上次的响应」；先返回错误也行，等结果落库了再升级。
var ErrDuplicateRequest = errors.New("duplicate request")
