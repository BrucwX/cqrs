package redis

import (
	"context"
	"time"
)

// Notify 生产延时任务
func (s *RedisData) Notify(ctx context.Context, productType string, orderID int64, delay time.Duration) error {
	zsetKey := s.getZSetKey(productType)
	streamKey := s.getStreamKey(productType)
	expireAt := time.Now().Add(delay).Unix()

	_, err := pushAndNotifyLua.Run(ctx, s.rdb,
		[]string{zsetKey, streamKey},
		expireAt, orderID,
	).Result()

	return err
}

// StartCloseCourse 专用于课程关单的消费启动器：
// 启动后台常驻协程，自动检测 ZSET，拉取已到期任务并交给 h 执行
func (s *RedisData) StartCloseCourse(ctx context.Context, h CourseHandler) {
	productType := "course"
	zsetKey := s.getZSetKey(productType)
	streamKey := s.getStreamKey(productType)
	groupName := s.getGroupName(productType)

	// 初始化消费组（不存在时自动创建）
	_ = s.rdb.XGroupCreateMkStream(ctx, streamKey, groupName, "0").Err()

	// pending 保存上一轮被唤醒、但尚未确认的唤醒信号。
	// 唤醒信号不能收到就 ACK：必须等本轮重新扫描任务、重新计算完睡眠时间后，
	// 它才真正“用过”，这时再清理；否则中间崩溃会丢掉这次唤醒。
	var pending []wakeupMsg

	// 启动常驻调度循环
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		// 1. 扫描并处理已到期的任务：先 handle，成功后才从 ZSET 移除
		// handle 出错时任务保留在 ZSET，等后续重试，不会丢单
		handleFailed := false
		for {
			// 只读取出到期任务，不删除
			taskIDs, err := s.peekDueTasks(ctx, zsetKey, 20)
			if err != nil || len(taskIDs) == 0 || h == nil {
				break
			}

			for _, id := range taskIDs {
				// 先执行 handle
				if err := h(ctx, id); err != nil {
					// 处理失败：不弹出，任务留在 ZSET 里等下一轮重试
					handleFailed = true
					break
				}
				// handle 成功，此时才弹出（从 ZSET 移除）
				if err := s.removeTask(ctx, zsetKey, id); err != nil {
					handleFailed = true
					break
				}
			}

			if handleFailed {
				break
			}
		}

		// 2. 检查未来最近任务的触发时间，重新计算本轮睡眠时间
		nextAt := s.getNextTaskTime(ctx, zsetKey)
		var sleepDuration time.Duration

		if nextAt == 0 {
			// 没有任务：进入深度睡眠，等待生产者 Notify 发信号唤醒
			sleepDuration = 0
		} else {
			// 有未来任务：计算时间差，睡到下一个任务的时间点
			sleepDuration = time.Until(time.Unix(nextAt, 0))
			if sleepDuration <= 0 {
				// 任务刚好或已到期，设为 -1 代表非阻塞立即进入下一轮拉取
				sleepDuration = -1
			}
		}

		// 有任务处理失败、被留在 ZSET 时，给一个最小退避（它始终处于“已到期”状态），
		// 避免失败任务导致下一轮立即重试而空转；退避结束后仍然会重新取到它
		if handleFailed && sleepDuration <= 0 {
			sleepDuration = time.Second
		}

		// 3. 睡眠时间已重新计算完毕，上一轮的唤醒信号至此才算用完，现在清理
		s.ackWakeups(ctx, productType, pending)
		pending = nil

		// 4. 阻塞休眠（利用 XReadGroup 的 Block）
		// 收到的唤醒信号先不 ACK，等下一轮重算完睡眠时间后在第 3 步清理
		pending = s.sleep(ctx, productType, streamKey, sleepDuration)
	}
}
