package job

import (
	"context"

	"go.uber.org/zap"
)

type AwardStockUpdater interface {
	UpdateAwardStock(ctx context.Context) (bool, error)
}

type UpdateAwardStockJob struct {
	updater AwardStockUpdater
	logger  *zap.Logger
}

func NewUpdateAwardStockJob(updater AwardStockUpdater, logger *zap.Logger) *UpdateAwardStockJob {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &UpdateAwardStockJob{updater: updater, logger: logger}
}

func (j *UpdateAwardStockJob) Exec() {
	updated, err := j.updater.UpdateAwardStock(context.Background())
	if err != nil {
		j.logger.Error("update award stock job failed", zap.Error(err))
		return
	}
	if updated {
		j.logger.Info("update award stock job completed")
	}
}
