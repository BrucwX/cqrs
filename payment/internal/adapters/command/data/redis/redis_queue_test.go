package redis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// 测试常量
const (
	testProductType = "course"
	testOrderID     = 1001
	testPodID1      = "test-pod-1"
	testPodID2      = "test-pod-2"

	testZSetKey   = "delay:zset:course"
	testStreamKey = "delay:stream:course"
)

// setupTestRedis 连接本地真实 Redis，清理测试数据，返回 redis client 和清理函数
func setupTestRedis(t *testing.T) (*redis.Client, func()) {
	t.Helper()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	ctx := context.Background()
	// 验证连接是否正常
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Fatalf("failed to connect to Redis: %v", err)
	}
	// 清理测试相关的 key
	// rdb.Del(ctx, testZSetKey, testStreamKey)
	return rdb, func() {
		//rdb.Del(ctx, testZSetKey, testStreamKey)
		rdb.Close()
	}
}

func HandlerTest(ctx context.Context, orderID int64) error {
	fmt.Printf("Processing order: %d\n", orderID)
	return nil
}

// testTask 定义测试任务的结构
type testTask struct {
	OrderID int64
	Score   float64
}

func createTestTask(ctx context.Context, rdb *redis.Client, zsetKey, streamKey string, tasks ...testTask) {
	for _, task := range tasks {
		rdb.ZAdd(ctx, zsetKey, redis.Z{Score: task.Score, Member: task.OrderID})
	}
	// 只需一条 wakeup 信号即可唤醒消费者
	rdb.XAdd(ctx, &redis.XAddArgs{Stream: streamKey, Values: []interface{}{"action", "wakeup"}})
}

func createTestTaskSilent(ctx context.Context, rdb *redis.Client, zsetKey, streamKey string, tasks ...testTask) {
	for _, task := range tasks {
		rdb.ZAdd(ctx, zsetKey, redis.Z{Score: task.Score, Member: task.OrderID})
	}
}

// 生成测试数据
func TestGenerateTask(t *testing.T) {
	rdb, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	_ = NewredisData(rdb, testPodID1)

	tasklist := []testTask{
		{OrderID: 1001, Score: 0},
		{OrderID: 1002, Score: float64(time.Now().Add(1 * time.Second).Unix())},
		{OrderID: 1003, Score: float64(time.Now().Add(3 * time.Second).Unix())},
		{OrderID: 1004, Score: float64(time.Now().Add(4 * time.Second).Unix())},
		{OrderID: 1005, Score: float64(time.Now().Add(5 * time.Second).Unix())},

		{OrderID: 1010, Score: float64(time.Now().Add(1 * time.Minute).Unix())},
		{OrderID: 1020, Score: float64(time.Now().Add(3 * time.Minute).Unix())},

		{OrderID: 1100, Score: float64(time.Now().Add(1 * time.Hour).Unix())},
	}

	// 用 Redis 原生命令直接写入多个已到期的任务（不发送 wakeup 信号）
	createTestTaskSilent(ctx, rdb, testZSetKey, testStreamKey, tasklist...)
}

func TestRestart(t *testing.T) {
	rdb, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	scheduler := NewredisData(rdb, testPodID1)

	handler := func(ctx context.Context, orderID int64) error {
		fmt.Printf("Processing order: %d\n", orderID)
		return nil
	}

	go scheduler.StartCloseCourse(ctx, handler)

	time.Sleep(1 * time.Second)
}

// 没有任务时,测试是否sleep (阻塞)
func TestSleepingWithContext(t *testing.T) {
	rdb, cleanup := setupTestRedis(t)
	defer cleanup()

	// 控制沉睡观测时间
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 确保测试结束时停止后台任务

	scheduler := NewredisData(rdb, testPodID1)

	done := make(chan struct{})

	go func() {
		scheduler.StartCloseCourse(ctx, HandlerTest)
		close(done)
	}()

	select {
	case <-done:
		t.Fatal("期望任务处于沉睡/阻塞状态，但它提前退出了")
	case <-time.After(10 * time.Second):
		// 100ms 内没有退出，确认当前处于沉睡状态，测试通过
	}
}

// 提前用 Redis 原生命令往里面写多个任务，然后再启动消费者
// 验证能否处理历史存量任务（模拟服务重启场景）
func TestRestartWithGenerateTask(t *testing.T) {
	rdb, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()
	scheduler := NewredisData(rdb, testPodID1)

	tasklist := []testTask{
		{OrderID: 1001, Score: 0},
		{OrderID: 1002, Score: 0},
		{OrderID: 1003, Score: 0},
		{OrderID: 1003, Score: 500000000},
	}

	// 用 Redis 原生命令直接写入多个已到期的任务（不发送 wakeup 信号）
	createTestTaskSilent(ctx, rdb, testZSetKey, testStreamKey, tasklist...)

	// 启动消费者，验证能否处理历史存量任务
	processed := make(chan int64, len(tasklist))
	handler := func(ctx context.Context, orderID int64) error {
		processed <- orderID
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go scheduler.StartCloseCourse(ctx, handler)

	// 从 channel 收集已处理的任务
	processedIDs := make(map[int64]bool)
	timeout := time.After(5 * time.Second)
	for range len(tasklist) {
		select {
		case id := <-processed:
			processedIDs[id] = true
		case <-timeout:
			t.Fatal("timeout waiting for tasks to be processed")
		}
	}

	// 验证所有任务都被处理了
	for _, task := range tasklist {
		if !processedIDs[task.OrderID] {
			t.Fatalf("expected orderID %d to be processed, but it wasn't", task.OrderID)
		}
	}

	// 验证 ZSET 已被清空
	count, err := rdb.ZCard(ctx, testZSetKey).Result()
	if err != nil {
		t.Fatalf("ZCard failed: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected ZSET to be empty after processing, got %d", count)
	}

}
