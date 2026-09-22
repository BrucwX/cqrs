package ports

import (
	"context"
	"errors"
	"sync"
	"time"
)

// memoryClosePayQueue 是 ClosePayQueue 的内存实现，**只给本包的单元测试用**（所以不导出）：
// 拿它把端口的语义钉死、能跑能测，真实实现（Redis 等）照着它的行为对齐。
//
// 语义跟接口注释一一对应：
//   - Notify 记下「下次该投的时刻」= now + delay（同一个 orderID 重复 Notify = 重设时刻）；
//   - Listen **阻塞式**：睡到最近一条到点，或被新入队唤醒后重算；一次最多取 n 条投到通道上
//     （缓冲 = n）；投出去没 Ack 的，过 redeliver 再投一遍（至少一次）；
//   - Ack 销掉一条；
//   - 通道谁建谁关：ctx 结束就关。
type memoryClosePayQueue struct {
	mu   sync.Mutex
	next map[int64]time.Time // orderID → 下次该投的时刻（Ack 之后从表里消失）

	redeliver time.Duration // 投出去多久没 Ack 就重投
	batch     int           // Start 一次最多取几条
	wake      chan struct{} // 有人入队时敲一下，让 Listen 重算该睡多久
}

var _ ClosePayQueue = (*memoryClosePayQueue)(nil)

// newMemoryClosePayQueue 建内存队列：redeliver = 投出去多久没 Ack 就重投，
// batch = Start 一次最多取几条。
func newMemoryClosePayQueue(redeliver time.Duration, batch int) *memoryClosePayQueue {
	if batch <= 0 {
		batch = 1 // 没给就一条一条来
	}

	return &memoryClosePayQueue{
		next:      make(map[int64]time.Time),
		redeliver: redeliver,
		batch:     batch,
		wake:      make(chan struct{}, 1),
	}
}

// handle 拿到信号怎么办事 —— 内存版没有真业务（真业务在用例里：CloseIfUnpaid），
// 就当「办成了」，好让 Start 把这条销掉。
func (q *memoryClosePayQueue) handle(_ context.Context, _ int64) error {
	return nil
}

func (q *memoryClosePayQueue) Notify(_ context.Context, orderID int64, delay time.Duration) error {
	q.mu.Lock()
	q.next[orderID] = time.Now().Add(delay)
	q.mu.Unlock()

	select {
	case q.wake <- struct{}{}: // 叫醒正在睡的 Listen 重算
	default: // 已经有一声没被取走就算了
	}

	return nil
}

// Listen 是接口那一面：一次最多取 n 条，投到返回的通道上（阻塞式；ctx 结束关通道）。
func (q *memoryClosePayQueue) Listen(ctx context.Context, n int) (<-chan int64, error) {
	if n <= 0 {
		return nil, errors.New("ports: listen batch must be positive")
	}

	out := make(chan int64, n) // 缓冲 = n：正好装下这一批

	go func() {
		defer close(out) // 谁建谁关

		for {
			ids, wait := q.due(n)

			if len(ids) > 0 {
				for _, id := range ids {
					select {
					case out <- id: // 投不进去就阻塞在这里 —— 背压
					case <-ctx.Done():
						return
					}
				}
				continue
			}

			if wait > 0 { // 有信号但还没到点：睡到最近那条（或被新入队唤醒）
				if !sleep(ctx, q.wake, wait) {
					return
				}
				continue
			}

			// 队列是空的：没有「到点」这回事，就**干等被唤醒** ——
			// 别拿个「睡一小时」的定时器凑数（那又变成轮询了）。
			select {
			case <-ctx.Done():
				return
			case <-q.wake:
			}
		}
	}()

	return out, nil
}

// Ack 销掉一条：办成了才调它。没 Ack 的过 redeliver 会被重新投出来（至少一次）。
func (q *memoryClosePayQueue) Ack(_ context.Context, orderID int64) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	delete(q.next, orderID)
	return nil
}

// Start 是队列自己的入口（跟 closePayScheduler.Start 一个形状）：用自己的 batch 跑 ——
// 从 Listen 拿通道，逐条交给 handle 办，**办成了才 Ack**；阻塞到 ctx 结束。
func (q *memoryClosePayQueue) Start(ctx context.Context) error {
	signals, err := q.Listen(ctx, q.batch)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case orderID, ok := <-signals:
			if !ok {
				return nil // 通道关了（Listen 的 ctx 结束）
			}

			if err := q.handle(ctx, orderID); err != nil {
				continue // 不 Ack：过 redeliver 会再投一遍
			}
			if err := q.Ack(ctx, orderID); err != nil {
				return err
			}
		}
	}
}

// due 取最多 n 条到点的，并把它们的「下次该投」推后 redeliver；
// 同时告诉调用方还没到点的话该睡多久（wait > 0 = 睡到最近一条到点；wait = 0 = 队列是空的）。
func (q *memoryClosePayQueue) due(n int) (ids []int64, wait time.Duration) {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := time.Now()
	var (
		more     bool      // 到点了但这一批装不下
		earliest time.Time // 最近一条还没到点的时刻
	)

	for id, at := range q.next {
		switch {
		case !at.After(now) && len(ids) < n:
			ids = append(ids, id)
			q.next[id] = now.Add(q.redeliver)
		case !at.After(now):
			more = true
		case earliest.IsZero() || at.Before(earliest):
			earliest = at
		}
	}

	switch {
	case more:
		wait = time.Millisecond // 还有到点的没取完：马上回来
	case !earliest.IsZero():
		wait = time.Until(earliest) // 有信号还没到点：睡到它到点
	}
	// 队列空的：wait 留在 0 —— 意思是「没到点这回事，干等唤醒」，调用方别设定时器。

	return ids, wait
}

// sleep 睡到时刻到、或被唤醒；ctx 结束返回 false。
func sleep(ctx context.Context, wake <-chan struct{}, d time.Duration) bool {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-wake:
		return true
	case <-timer.C:
		return true
	}
}
