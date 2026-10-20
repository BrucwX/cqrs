package redis

import "fmt"

// Key 辅助函数
func (s *redisData) getZSetKey(productType string) string {
	return fmt.Sprintf("delay:zset:%s", productType)
}

func (s *redisData) getStreamKey(productType string) string {
	return fmt.Sprintf("delay:stream:%s", productType)
}

// getGroupName 消费者组名：和 key 一样只由 productType 推导，
// 不依赖外部注入，不同业务类型的消费进度天然隔离。
func (s *redisData) getGroupName(productType string) string {
	return fmt.Sprintf("delay:group:%s", productType)
}

// getConsumerName 消费者名：由 productType 和 podID 共同推导。
// podID 用于区分不同实例，支持多副本并行消费。
func (s *redisData) getConsumerName(productType string) string {
	return fmt.Sprintf("delay:consumer:%s:%s", productType, s.podID)
}
