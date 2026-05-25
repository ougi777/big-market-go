package repository

import (
	"context"

	"bm-go/internal/infrastructure/persistent/po"

	"gorm.io/gorm"
)

func (r *StrategyRepository) AwardStockConsumeSendQueue(ctx context.Context, strategyID int64, awardID int) error {
	if r.stockQueue != nil {
		return r.stockQueue.AwardStockConsumeSendQueue(ctx, strategyID, awardID)
	}
	return nil
}

func (r *StrategyRepository) UpdateStrategyAwardStock(ctx context.Context, strategyID int64, awardID int) error {
	return r.defaultDB(ctx).
		Model(&po.StrategyAward{}).
		Where("strategy_id = ? and award_id = ? and award_count_surplus > 0", strategyID, awardID).
		UpdateColumn("award_count_surplus", gorm.Expr("award_count_surplus - ?", 1)).
		Error
}
