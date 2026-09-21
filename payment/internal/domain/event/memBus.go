package event

import (
	"context"
	"errors"
	"sync"
)

// bus 是进程内事件总线的**实现**：订阅者按 topic 攒着，发布时逐个调用。
//
// 它是私有的 —— 外面只通过 EventBus 接口拿它（NewBus 返回接口），拿不到也不能换实现细节。
//
// 语义（也是单测钉住的那几条）：
//   - 同一个 topic 的订阅者**按订阅顺序**依次调用；
//   - Publish 同步跑完该 topic 的**所有**订阅者才返回，**一个订阅者失败不影响别人**，
//     错误汇总（errors.Join）返回；
//   - PublishAsync 起协程分发、立刻返回，错误交给构造时给的 onError。
//
// ⚠️ 它是**进程内**的：进程重启就没了，多实例各自广播（不是真·跨进程事件总线）。
// 要跨进程就另写一个实现（kafka / rocketmq / 本地消息表 + 定时投递），
// 两边只共享 EventBus 这个接口。
type bus struct {
	mu       sync.RWMutex
	handlers map[EventTopic][]EventHandler

	onError func(error)
}

// 编译期断言：实现必须满足事件总线端口。
var _ EventBus = (*bus)(nil)

// NewBus 构造进程内总线，**返回接口**（具体实现私有，外面只认 EventBus）。
//
// onError 接收**异步分发**时订阅者返回的错误（同步 Publish 直接把错误返回给调用方）；
// 传 nil 就丢弃（异步失败总得有个去处，生产上建议传日志 / 告警）。
func NewBus(onError func(error)) EventBus {
	if onError == nil {
		onError = func(error) {}
	}
	return &bus{
		handlers: make(map[EventTopic][]EventHandler),
		onError:  onError,
	}
}

// Subscribe 订阅某个 topic。
//
// 同一个 handler 订阅多次就会被调用多次（不做去重）；nil handler 直接忽略。
func (b *bus) Subscribe(topic EventTopic, handler EventHandler) {
	if handler == nil {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[topic] = append(b.handlers[topic], handler)
}

// Publish 同步广播：跑完该 topic 的所有订阅者才返回，失败汇总返回（没失败就 nil）。
func (b *bus) Publish(ctx context.Context, e Event) error {
	if e == nil {
		return nil
	}

	var errs []error
	for _, h := range b.subscribers(e.Topic()) {
		if err := h(ctx, e); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// PublishAsync 异步广播：另起协程分发，立刻返回；失败交给 onError。
func (b *bus) PublishAsync(ctx context.Context, e Event) {
	go func() {
		if err := b.Publish(ctx, e); err != nil {
			b.onError(err)
		}
	}()
}

// subscribers 取某个 topic 的订阅者快照（复制一份，别让发布期间的新订阅改到正在遍历的切片）。
func (b *bus) subscribers(topic EventTopic) []EventHandler {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return append([]EventHandler(nil), b.handlers[topic]...)
}
