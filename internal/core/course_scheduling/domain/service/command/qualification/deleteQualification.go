package qualification

import "context"

// DeleteQualification 删除授课资质
//
// 资质不存在时报 ErrQualificationNotFound。
func (h *Handler) DeleteQualification(ctx context.Context, id int64) error {
	return h.QualificationCmd.Delete(id)
}
