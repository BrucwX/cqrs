package qualification

import "context"

// DeleteQualification 删除授课资质
//
// 资质不存在时报 ErrQualificationNotFound。
func (h *Handler) DeleteQualification(ctx context.Context, id int64) (err error) {
	ctx, err = h.tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { err = h.tx.End(ctx, err) }()

	return h.QualificationCmd.DeleteQualification(ctx, id)
}
