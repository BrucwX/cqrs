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
}

// PayIntent 一次支付的「意图」：谁、哪张单、多少钱、走哪条路。
//
// 渠道专属的东西（openid / return_url / 场景）不进领域，丢进 Scene / Extra，
// 由渠道实现自己解读。
type PayIntent struct {
	OrderID    int64
	CustomerID int64
	Amount     int64            // 分：跟订单 payableAmount 对齐
	Channel    order.PayChannel // 支付途径（微信 / 支付宝 / 银行卡…）
	Scene      string           // 渠道内场景：app / h5 / native / jsapi
	Extra      map[string]string
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
