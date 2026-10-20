package mysql

import (
	"context"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v3/log"

	appcmd "courseScheduling/internal/app/command/appRepo"
	"courseScheduling/internal/domain/aggregate/classroom"
)

const maxOptimisticRetries = 10
const retryPause = 5 * time.Millisecond

func (d *MysqlData) RunUpdatClassroom(
	ctx context.Context,
	id string,
	updateFn func(ctx context.Context, cl *classroom.Classroom) (*classroom.Classroom, error),
) (err error) {
	// 一阶段：乐观重试。不加锁，输了就重读最新版本再来。
	for attempt := 1; attempt <= maxOptimisticRetries; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		// 读教室（PO 里带着这一轮的 lock_version），把它交给 updateFn 改完写回
		po, err := d.classroomPO(ctx, id, false)
		if err != nil {
			return err
		}

		err = d.updateClassroomPO(ctx, po, updateFn)
		if !errors.Is(err, appcmd.ErrVersionConflict) {
			return err
		}

		// 冲突了（最后一次也照实说「冲突」—— 紧接着就去二阶段了）
		log.Info("optimistic lock: version conflict", "attempt", attempt, "version", po.LockVersion)
		if !pauseOrDone(ctx, retryPause) {
			return ctx.Err()
		}
	}

	// 二阶段：正经事务。读带 FOR UPDATE 锁住行，这一轮别人插不进来，必然写成。
	log.Info("optimistic lock: retries exhausted, falling back to a real transaction",
		"attempts", maxOptimisticRetries)

	tctx, err := d.Begin(ctx)
	if err != nil {
		return err
	}
	// 出错回滚（正常路径走下面的 commit）；ctx 里的事务在 defer 执行时还没结束。
	defer func() {
		if err != nil {
			if rbErr := d.rollback(tctx); rbErr != nil {
				log.Error("mysql rollback", "err", rbErr)
			}
		}
	}()

	err = d.UpdateClassroom(tctx, id, updateFn)
	if err != nil {
		return err
	}
	err = d.commit(tctx)
	return
}

// pauseOrDone 停一小会儿；ctx 结束返回 false。
func pauseOrDone(ctx context.Context, d time.Duration) bool {
	if d <= 0 {
		return ctx.Err() == nil
	}

	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}
