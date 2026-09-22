package ports

import (
	"context"
	"testing"
	"time"
)

// produce 往 q 里塞 n 条信号（orderID = 1..n，立刻到点）——
// 测试用的帮手，**不是测试函数**：n 就是「这次让队列产生几条信号」。
func produce(q *memoryClosePayQueue, n int) {
	for i := 1; i <= n; i++ {
		// 内存实现不会失败（没网络、没容量限制），真失败就是参数写错了。
		_ = q.Notify(context.Background(), int64(i), 0)
	}
}

// 消费：Listen 一次最多给 n 条，多出来的接着给 —— 这里 5 条信号、一次只取 2 条。
func TestMemoryQueueListen(t *testing.T) {
	const total = 5

	q := newMemoryClosePayQueue(50*time.Millisecond, 2)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	produce(q, total) // 往队列里塞 5 条（立刻到点）

	signals, err := q.Listen(ctx, 2) // 一次最多取 2 条
	if err != nil {
		t.Fatalf("Listen: %v", err)
	}

	// 5 条一条都不能少（队列会分 2 + 2 + 1 三轮把它们取出来）。
	seen := make(map[int64]bool, total)
	for i := 0; i < total; i++ {
		select {
		case id := <-signals:
			seen[id] = true
			if err := q.Ack(ctx, id); err != nil { // 办完了就销掉
				t.Fatalf("Ack: %v", err)
			}
		case <-time.After(time.Second):
			t.Fatalf("只收到 %d 条，想要 %d 条", i, total)
		}
	}

	if len(seen) != total {
		t.Fatalf("收到 %d 个不同的信号，想要 %d 个", len(seen), total)
	}
}
