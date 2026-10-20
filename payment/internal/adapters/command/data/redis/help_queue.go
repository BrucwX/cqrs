package redis

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// peekDueTasks 读取已到期任务 ID（只读，不删除）。
// 消费者只有一个，不存在并发重复拉取，所以不需要 Lua 原子弹出；
// 任务在 handle 成功返回后才由 removeTask 从 ZSET 中移除，
// 这样 handle 出错或进程崩溃时任务仍留在 ZSET，下一轮会重新取到，不会丢单。
func (s *redisData) peekDueTasks(ctx context.Context, zsetKey string, limit int) ([]int64, error) {
	// 用 ZRangeArgs + ByScore 替代已废弃的 ZRangeByScore（Redis 6.2.0+），语义等价：
	// ZRANGE key 0 <now> BYSCORE LIMIT 0 <limit>，start/stop 均为闭区间
	now := time.Now().Unix()
	vals, err := s.rdb.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:     zsetKey,
		Start:   0,
		Stop:    now,
		ByScore: true,
		Count:   int64(limit),
	}).Result()
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(vals))
	for _, val := range vals {
		if id, err := strconv.ParseInt(val, 10, 64); err == nil {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

// removeTask 任务处理成功后，才把它从 ZSET 中弹出
func (s *redisData) removeTask(ctx context.Context, zsetKey string, orderID int64) error {
	return s.rdb.ZRem(ctx, zsetKey, orderID).Err()
}

// getNextTaskTime 查询最早到期的任务时间戳
func (s *redisData) getNextTaskTime(ctx context.Context, zsetKey string) int64 {
	tasks, err := s.rdb.ZRangeWithScores(ctx, zsetKey, 0, 0).Result()
	if err != nil || len(tasks) == 0 {
		return 0
	}
	return int64(tasks[0].Score)
}

// wakeupMsg 一条已读取但尚未确认的唤醒信号
// 记录所属 stream，便于后续按原 key 精确 ACK
type wakeupMsg struct {
	stream string
	id     string
}

// sleep 挂起等待唤醒，返回本次收到的唤醒信号（此处不做 ACK）。
// 调用方必须在「重新扫描到期任务 + 重新计算睡眠时间」完成之后，
// 再用 ackWakeups 清理这些信号，保证唤醒不会因为中途崩溃而丢失。
func (s *redisData) sleep(ctx context.Context, productType, streamKey string, duration time.Duration) []wakeupMsg {
	groupName := s.getGroupName(productType)

	args := &redis.XReadGroupArgs{
		Group:    groupName,
		Consumer: s.getConsumerName(productType),
		Streams:  []string{streamKey, ">"},
		Count:    10,
	}
	if duration >= 0 {
		args.Block = duration
	}

	streams, err := s.rdb.XReadGroup(ctx, args).Result()
	if errors.Is(err, redis.Nil) || err != nil {
		return nil
	}

	msgs := make([]wakeupMsg, 0, len(streams))
	for _, stream := range streams {
		for _, msg := range stream.Messages {
			msgs = append(msgs, wakeupMsg{stream: stream.Stream, id: msg.ID})
		}
	}
	return msgs
}

// ackWakeups 清理唤醒信号。
// 只能在重新计算睡眠时间（以及处理完到期任务）之后调用，
// 此时信号已完成使命，提前确认会导致崩溃时丢唤醒。
func (s *redisData) ackWakeups(ctx context.Context, productType string, msgs []wakeupMsg) {
	if len(msgs) == 0 {
		return
	}
	groupName := s.getGroupName(productType)
	for _, msg := range msgs {
		s.rdb.XAck(ctx, msg.stream, groupName, msg.id)
	}
}
