package ports

import "context"

// 可能多个人有关闭支付的动作 , 所以这个是多对一消费
type ClosePayProducer interface {
	// 生产者信号
	ClosePayNotify(ctx context.Context, orderID int64) error
}

type ClosePayConsumer interface {
	// 定时轮询
	Start(ctx context.Context) error
	// 具体的消费行为函数
	CloseIfUnpaid(ctx context.Context, orderID int64) error
}
