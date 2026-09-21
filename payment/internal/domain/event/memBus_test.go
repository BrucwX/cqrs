package event

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// newEvent 造一个领域事件。事件接口是密封的，只有本包能自己造事件 ——
// 所以这个文件里可以顺手定义第二个事件，用来证明「按 topic 分发」。
func newEvent() *OrderPlacedEvent {
	return &OrderPlacedEvent{
		EventMeta:  NewEventMeta("1001"),
		CustomerId: "3001",
	}
}

// testEvent 测试用的第二个事件（自有 topic）。
type testEvent struct {
	EventMeta
}

func (testEvent) Topic() EventTopic { return "test.another" }

func (testEvent) isDomainEvent() {}

// recorder 按顺序记下被调用时收到的事件，返回 nil。
func recorder(hits *[]Event) EventHandler {
	return func(_ context.Context, e Event) error {
		*hits = append(*hits, e)
		return nil
	}
}

func TestPublish(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("订阅了的 topic 收到事件", func(t *testing.T) {
		t.Parallel()

		bus := NewBus(nil)
		e := newEvent()

		var hits []Event
		bus.Subscribe(e.Topic(), recorder(&hits))

		if err := bus.Publish(ctx, e); err != nil {
			t.Fatalf("不该有错误，实际 %v", err)
		}
		if len(hits) != 1 {
			t.Fatalf("订阅者该被调用 1 次，实际 %d 次", len(hits))
		}
		if hits[0] != Event(e) {
			t.Fatal("收到的不是同一个事件")
		}
	})

	t.Run("多个订阅者按订阅顺序依次调用", func(t *testing.T) {
		t.Parallel()

		bus := NewBus(nil)
		e := newEvent()

		var order []string
		bus.Subscribe(e.Topic(), func(context.Context, Event) error {
			order = append(order, "第一")
			return nil
		})
		bus.Subscribe(e.Topic(), func(context.Context, Event) error {
			order = append(order, "第二")
			return nil
		})
		bus.Subscribe(e.Topic(), func(context.Context, Event) error {
			order = append(order, "第三")
			return nil
		})

		if err := bus.Publish(ctx, e); err != nil {
			t.Fatalf("不该有错误，实际 %v", err)
		}
		if len(order) != 3 || order[0] != "第一" || order[1] != "第二" || order[2] != "第三" {
			t.Fatalf("调用顺序不对：%v", order)
		}
	})

	t.Run("按 topic 分发：别的 topic 的订阅者不受影响", func(t *testing.T) {
		t.Parallel()

		bus := NewBus(nil)
		placed := newEvent()
		another := testEvent{EventMeta: NewEventMeta("2002")}

		var placedHits, anotherHits []Event
		bus.Subscribe(placed.Topic(), recorder(&placedHits))
		bus.Subscribe(another.Topic(), recorder(&anotherHits))

		if err := bus.Publish(ctx, placed); err != nil {
			t.Fatalf("不该有错误，实际 %v", err)
		}
		if len(placedHits) != 1 || len(anotherHits) != 0 {
			t.Fatalf("只该本 topic 的订阅者被调：placed=%d another=%d", len(placedHits), len(anotherHits))
		}

		if err := bus.Publish(ctx, another); err != nil {
			t.Fatalf("不该有错误，实际 %v", err)
		}
		if len(anotherHits) != 1 || len(placedHits) != 1 {
			t.Fatalf("另一个 topic 也该各收各的：placed=%d another=%d", len(placedHits), len(anotherHits))
		}
	})

	t.Run("没有订阅者 → 不算错", func(t *testing.T) {
		t.Parallel()

		bus := NewBus(nil)

		if err := bus.Publish(ctx, newEvent()); err != nil {
			t.Fatalf("没有订阅者时不该报错，实际 %v", err)
		}
	})

	t.Run("空事件 → 直接返回", func(t *testing.T) {
		t.Parallel()

		bus := NewBus(nil)

		if err := bus.Publish(ctx, nil); err != nil {
			t.Fatalf("空事件不该报错，实际 %v", err)
		}
	})

	t.Run("nil handler 会被忽略", func(t *testing.T) {
		t.Parallel()

		bus := NewBus(nil)
		e := newEvent()

		var hits []Event
		bus.Subscribe(e.Topic(), nil)
		bus.Subscribe(e.Topic(), recorder(&hits))

		if err := bus.Publish(ctx, e); err != nil {
			t.Fatalf("不该有错误，实际 %v", err)
		}
		if len(hits) != 1 {
			t.Fatalf("nil handler 该被忽略，订上的那个该被调 1 次，实际 %d", len(hits))
		}
	})

	t.Run("有订阅者报错 → 错误返回，但别的订阅者照跑", func(t *testing.T) {
		t.Parallel()

		boom := errors.New("订阅者一挂了")
		bus := NewBus(nil)
		e := newEvent()

		var hits []Event
		bus.Subscribe(e.Topic(), func(context.Context, Event) error { return boom })
		bus.Subscribe(e.Topic(), recorder(&hits))

		err := bus.Publish(ctx, e)
		if !errors.Is(err, boom) {
			t.Fatalf("期望错误 %v，实际 %v", boom, err)
		}
		if len(hits) != 1 {
			t.Fatalf("前一个订阅者失败不该挡住后面的，实际只跑了 %d 个", len(hits))
		}
	})

	t.Run("多个订阅者报错 → 全部汇总", func(t *testing.T) {
		t.Parallel()

		first := errors.New("第一个错")
		second := errors.New("第二个错")
		bus := NewBus(nil)
		e := newEvent()

		bus.Subscribe(e.Topic(), func(context.Context, Event) error { return first })
		bus.Subscribe(e.Topic(), func(context.Context, Event) error { return second })

		err := bus.Publish(ctx, e)
		if !errors.Is(err, first) || !errors.Is(err, second) {
			t.Fatalf("两个错误都该在返回值里，实际 %v", err)
		}
	})

	t.Run("同一个 handler 订两次会被调两次", func(t *testing.T) {
		t.Parallel()

		bus := NewBus(nil)
		e := newEvent()

		var hits []Event
		h := recorder(&hits)
		bus.Subscribe(e.Topic(), h)
		bus.Subscribe(e.Topic(), h)

		if err := bus.Publish(ctx, e); err != nil {
			t.Fatalf("不该有错误，实际 %v", err)
		}
		if len(hits) != 2 {
			t.Fatalf("订两次该调两次，实际 %d 次", len(hits))
		}
	})
}

func TestPublishAsync(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("不等订阅者跑完就返回", func(t *testing.T) {
		t.Parallel()

		bus := NewBus(nil)
		e := newEvent()

		block := make(chan struct{})
		started := make(chan struct{})
		bus.Subscribe(e.Topic(), func(context.Context, Event) error {
			close(started)
			<-block // 卡住订阅者
			return nil
		})

		returned := make(chan struct{})
		go func() {
			bus.PublishAsync(ctx, e)
			close(returned)
		}()

		select {
		case <-returned:
			// 预期：分发是异步的，不等订阅者
		case <-time.After(time.Second):
			t.Fatal("PublishAsync 不该等订阅者跑完")
		}

		select {
		case <-started:
		case <-time.After(time.Second):
			t.Fatal("订阅者该在协程里被调到")
		}
		close(block)
	})

	t.Run("订阅者报错 → 交给 onError", func(t *testing.T) {
		t.Parallel()

		boom := errors.New("异步订阅者挂了")
		errCh := make(chan error, 1)
		bus := NewBus(func(err error) { errCh <- err })
		e := newEvent()

		bus.Subscribe(e.Topic(), func(context.Context, Event) error { return boom })
		bus.PublishAsync(ctx, e)

		select {
		case err := <-errCh:
			if !errors.Is(err, boom) {
				t.Fatalf("期望错误 %v，实际 %v", boom, err)
			}
		case <-time.After(time.Second):
			t.Fatal("异步失败该走 onError")
		}
	})

	t.Run("没订阅者也不报错", func(t *testing.T) {
		t.Parallel()

		errCh := make(chan error, 1)
		bus := NewBus(func(err error) { errCh <- err })

		bus.PublishAsync(ctx, newEvent())

		select {
		case err := <-errCh:
			t.Fatalf("没订阅者不该报错，实际 %v", err)
		case <-time.After(50 * time.Millisecond):
			// 预期：什么都不发生
		}
	})
}

// TestBusConcurrent 并发订阅 + 发布不炸（配合 -race 跑）。
func TestBusConcurrent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	bus := NewBus(nil)
	e := newEvent()

	var calls atomic.Int64

	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 20 {
				bus.Subscribe(e.Topic(), func(context.Context, Event) error {
					calls.Add(1)
					return nil
				})
				_ = bus.Publish(ctx, e)
			}
		}()
	}
	wg.Wait()

	if calls.Load() == 0 {
		t.Fatal("订阅者一次都没被调到")
	}
}
