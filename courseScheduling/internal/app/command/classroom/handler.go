package classroom

import (
	appcommand "courseScheduling/internal/app/command/appRepo"
	"courseScheduling/internal/domain/repo/command"
)

// Handler 教室命令处理器
type Handler struct {
	ClassroomCmd command.ClassroomCommand
	// tx 真事务：DeleteClassroom 这类「一把成」的用例用它。
	tx appcommand.Transaction
	// optTx 乐观锁事务：UpdateClassroom 用它 —— 读版本、写时校验、冲突重跑。
	optTx appcommand.OptimisticTransaction
}

// NewHandler 创建教室命令处理器
func NewHandler(
	cc command.ClassroomCommand,
	tx appcommand.Transaction,
	optTx appcommand.OptimisticTransaction,
) *Handler {
	return &Handler{
		ClassroomCmd: cc,
		tx:           tx,
		optTx:        optTx,
	}
}
