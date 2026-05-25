package job

import (
	"context"

	"go.uber.org/zap"
)

type ActivitySkuStockUpdater interface {
	UpdateActivitySkuStock(ctx context.Context) (bool, error)
}

type UpdateActivitySkuStockJob struct {
	updater ActivitySkuStockUpdater
	logger  *zap.Logger
}

func NewUpdateActivitySkuStockJob(updater ActivitySkuStockUpdater, logger *zap.Logger) *UpdateActivitySkuStockJob {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &UpdateActivitySkuStockJob{updater: updater, logger: logger}
}

func (j *UpdateActivitySkuStockJob) Exec() {
	updated, err := j.updater.UpdateActivitySkuStock(context.Background())
	if err != nil {
		j.logger.Error("update activity sku stock job failed", zap.Error(err))
		return
	}
	if updated {
		j.logger.Info("update activity sku stock job completed")
	}
}
