package order

import "strings"

// --- 值对象 (Value Object) ---

// PayChannel 是「支付途径」值对象：这笔钱是从**哪条路**进来的（微信 / 支付宝 / 余额 / 银行卡 / 现金…）。
//
// 值对象没有自己的 id，跟 product.Type 一个路子：
//   - **code 就是它的标识**，也是别处引用、对账、回调路由用的键（约定小写下划线：wechat / alipay /
//     balance / bank_card / cash），建好后不该改
//   - name 只是展示名（"微信支付" / "余额支付"），可以改，不影响任何判断
//
// ⚠️ 这里只记「走哪条路」，**流水号与金额在 PaidOrder 上**（tradeNo / paidAmount），
// 更细的支付流水（回调、渠道单号、手续费）属于 payment 那边的支付聚合，不往订单里塞。
type PayChannel struct {
	code string // 稳定键：wechat / alipay / balance / bank_card / cash
	name string // 展示名
}

// NewPayChannel 构造支付途径值对象：code 与 name 都必填。
func NewPayChannel(code string, name string) (PayChannel, error) {
	c := PayChannel{code: code, name: name}
	if !c.valid() {
		return PayChannel{}, ErrInvalidPayChannel
	}
	return c, nil
}

// ReconstitutePayChannel 仓储恢复：字段全部由外部给定，不做校验。
func ReconstitutePayChannel(code string, name string) PayChannel {
	return PayChannel{code: code, name: name}
}

// valid 是否是构造出来的合法途径（挡住零值 PayChannel）。
func (c PayChannel) valid() bool {
	return strings.TrimSpace(c.code) != "" && strings.TrimSpace(c.name) != ""
}

// Is 两条途径是不是同一条（按 code 判）。
func (c PayChannel) Is(other PayChannel) bool { return c.code == other.code }

// --- 属性只读访问器 (Getters) ---

func (c PayChannel) Code() string { return c.code }
func (c PayChannel) Name() string { return c.name }
