package appRepo

import (
	"context"

	"payment/internal/domain/aggregate/order"
)

// PayService 支付端口：应用层视角、渠道无关。
//
// 「发起支付」与「确认结果」分开：有些渠道发起后要等用户扫码 / 唤起，结果靠回调或查单，
// 而不是同步返回；少数同步渠道（余额、部分快捷）会直接给结果。
//
// ⚠️ 它不是某个聚合的存储接口：实现要调支付渠道 SDK，渠道的验签、证书、报文差异全封在实现里。
type PayService interface {
	// Prepare 发起支付：返回「前端下一步怎么办」（二维码 / 跳转链接 / 预支付参数），
	// 少数同步渠道会直接返回 success。
	Prepare(ctx context.Context, intent PayIntent) (PayResult, error)

	// Query 查一笔支付在渠道侧的状态：没有回调的渠道、对账、超时补偿用。
	Query(ctx context.Context, orderID int64) (PayStatus, error)

	// Close 关闭渠道侧订单：关单 / 超时后，让渠道那边也把单关掉
	// （微信有 closeorder，支付宝超时自动关，调了也无妨）。渠道侧关成功才返回 nil。
	Close(ctx context.Context, orderID int64) error

	// Refund 申请退款：让渠道把支付的钱**原路退回**。
	// 返回渠道的退款流水号与受理结果；异步打款的渠道，最终结果还要靠退款回调 / 查单确认。
	Refund(ctx context.Context, req RefundRequest) (RefundResult, error)
}

// PayIntent 一次支付的「意图」：谁、哪张单、多少钱、走哪条路。
//
// 标识分两层，别混：
//   - OrderID / CustomerID 是**我方**的主键（商户单号就用订单号，见 order.MarkPaid）；
//   - ThirdPartyOrderID / ThirdPartyCustomerID 是**第三方平台侧**的对应标识，
//     查单、对账、退款关联都得靠它们。外部 id 一律 string：openid、平台单号
//     的形态由平台定，别假设它一定是数字。
//
// 其余渠道专属的东西（return_url / 商品描述 / 证书序列号…）仍丢进 Scene / Extra，
// 由渠道实现自己解读。
type PayIntent struct {
	OrderID              int64  // 我方订单号（也就是交给平台的商户单号）
	CustomerID           int64  // 我方顾客 id
	ThirdPartyOrderID    string // 第三方平台的订单 id：这笔单在平台侧的唯一标识
	ThirdPartyCustomerID string // 第三方平台的顾客 id：微信 openid / 支付宝 buyer_id 等

	Amount         int64            // 分：跟订单 payableAmount 对齐
	Channel        order.PayChannel // 支付途径（微信 / 支付宝 / 银行卡…）
	Scene          string           // 渠道内场景：app / h5 / native / jsapi
	IdempotencyKey string           // 等幂key
	Extra          map[string]string
}

// PayResult 一次 Prepare 的结果：前端下一步怎么办，或同步渠道的最终结果。
type PayResult struct {
	Status  PayStatus // 等用户付 / 同步成功 / 失败 / 已关单
	PayInfo string    // 二维码 URL、跳转链接或预支付参数（给前端）
	TradeNo string    // 渠道流水号（同步成功时才有意义）
}

// PayStatus 一笔支付在渠道侧的状态。
type PayStatus int

const (
	PayPending PayStatus = iota + 1 // 已发起，等用户付款（异步渠道常态）
	PaySuccess                      // 已成功
	PayFailed                       // 失败
	PayClosed                       // 已关单 / 超时
)

// RefundRequest 一次退款申请：哪笔订单、退多少钱、退款单号。
//
// 退款是「原路退回」，所以要带上原支付的渠道流水号；退款单号（商户侧）我们自己生成。
type RefundRequest struct {
	OrderID  int64  // 原订单（也当商户单号用）
	RefundID int64  // 商户退款单号（refund 聚合的 id）
	TradeNo  string // 原支付的渠道流水号
	Amount   int64  // 退款金额（分）
	Reason   string // 退款原因
}

// RefundResult 一次退款申请的结果。
type RefundResult struct {
	RefundTradeNo string       // 渠道退款流水号
	Status        RefundStatus // 已受理 / 已退 / 失败
}

// RefundStatus 退款在渠道侧的状态。
type RefundStatus int

const (
	RefundPending RefundStatus = iota + 1 // 已受理，等渠道打款（异步渠道常态）
	RefundSuccess                         // 已退
	RefundFailed                          // 退款失败
)
