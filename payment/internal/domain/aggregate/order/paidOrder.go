package order

import "time"

// --- 已支付 ---

// PaidOrder 已支付状态：钱收到了，订单到这就停住（终态）。
type PaidOrder struct {
	OrderMeta
	paidAmount int64
	paidAt     time.Time
	tradeNo    string     // 支付流水号（第三方平台交易号），对账用
	payChannel PayChannel // 支付途径：钱从哪条路进来的（微信 / 支付宝 / 余额…）
}

// ReconstitutePaid 仓储恢复：字段全部由外部给定，不做校验。
func ReconstitutePaid(
	meta OrderMeta,
	paidAmount int64,
	paidAt time.Time,
	tradeNo string,
	payChannel PayChannel,
) PaidOrder {
	return PaidOrder{
		OrderMeta:  meta,
		paidAmount: paidAmount,
		paidAt:     paidAt,
		tradeNo:    tradeNo,
		payChannel: payChannel,
	}
}

func (o PaidOrder) PaidAmount() int64      { return o.paidAmount }
func (o PaidOrder) PaidAt() time.Time      { return o.paidAt }
func (o PaidOrder) TradeNo() string        { return o.tradeNo }
func (o PaidOrder) PayChannel() PayChannel { return o.payChannel }

// 已支付是终态：订单这边不再往外迁移 —— 退款走 refund 聚合（退多少钱、退几件、
// 审批与打款都在那边），订单只负责「钱收到了」这件事。
