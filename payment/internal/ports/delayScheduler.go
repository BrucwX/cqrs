package ports

import (
	"context"
	"time"
)

type CourseHandler func(ctx context.Context, orderID int64) error

type DelayScheduler interface {
	Notify(ctx context.Context, productType string, orderID int64, delay time.Duration) error

	StartCloseCourse(ctx context.Context, h CourseHandler)
}
