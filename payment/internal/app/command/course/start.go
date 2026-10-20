package course

import "context"

// Start / Stop：让 Handler 顺便当「一个 server」用。
//
// kratos.Server(...) 只要形状对上 —— Start(ctx) error + Stop(ctx) error ——
// 不用 import transport。ctx 也正来自这里：kratos.New() 内部有自己的 ctx，
// App.Run() 把它派生成 octx 交给每个 server 的 Start(octx)；进程收到
// SIGTERM/SIGINT 时会 App.Stop() → cancel 掉它，再并发调 Stop()。
// 也就是说 **ctx 归进程生命周期管**，用例层只负责「给我一个 ctx，我就一直转到它结束」。
//
//
//	kratos.Server(servers.HTTP, servers.GRPC, handler)
//
// Start 阻塞到 ctx 结束：等信号、睡到最近一条到点、重投都是队列端口
// （ports.DelayScheduler）的事，这里只把「办一条」的动作交出去。单条失败不让
// 循环停 —— 队列是「至少一次」，失败那条之后会重新投递，重复消费由
// CloseCourseIfUnpaid 的幂等兜住。
func (h *Handler) Start(ctx context.Context) error {
	go h.queue.StartCloseCourse(ctx, h.CloseCourseIfUnpaid)

	return nil
}
