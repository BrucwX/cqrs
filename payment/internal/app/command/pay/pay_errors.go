package pay

import "errors"

// ErrUnexpectedNotify 已关单 / 已取消的订单却收到了「支付成功」回调。
//
// 钱收了、单没了，属于资金差异：交给人工 / 对账处理，不能当作幂等成功吞掉。
var ErrUnexpectedNotify = errors.New("payment notify for a non-pending order")
