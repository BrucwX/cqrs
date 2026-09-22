package ports

import (
	"context"
	"time"
)

type ClosePayQueue interface {
	Notify(ctx context.Context, orderID int64, delay time.Duration) error

	Listen(ctx context.Context, n int) (<-chan int64, error)

	Ack(ctx context.Context, orderID int64) error
}
